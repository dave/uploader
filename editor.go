package uploader

// notest

import (
	"context"
	"encoding/json"

	"net/http"

	"io/ioutil"

	"frizz.io/editor/client/actions"
	"frizz.io/editor/client/editable"
	"frizz.io/editor/client/editors"
	"frizz.io/editor/client/models"
	"frizz.io/editor/client/stores"
	"frizz.io/editor/client/views"
	"frizz.io/flux"
	"frizz.io/system"
	"frizz.io/system/node"
	"github.com/dave/kerr"
	"github.com/dave/vecty"
	"github.com/dave/vecty/elem"
	"github.com/dave/vecty/event"
	"github.com/dave/vecty/prop"
	"github.com/dave/vecty/style"
)

var _ editable.Editable = (*Imgur)(nil)

func (s *Imgur) Format(rule *system.RuleWrapper) editable.Format {
	return editable.Block
}

func (s *Imgur) EditorView(ctx context.Context, node *node.Node, format editable.Format) vecty.Component {
	return NewIconEditorView(ctx, node, format)
}

type iconEditorDataKeys int

const (
	iconEditorUploading iconEditorDataKeys = 1
)

type IconEditorView struct {
	*views.View

	model  *models.EditorModel
	url    *models.EditorModel
	imgur  *Imgur
	editor *editors.StringEditorView
}

func NewIconEditorView(ctx context.Context, node *node.Node, format editable.Format) *IconEditorView {
	v := &IconEditorView{}
	v.View = views.New(ctx, v)
	v.model = v.App.Editors.Get(node)
	v.url = v.App.Editors.Get(node.Map["url"])
	v.imgur = v.model.Node.Value.(*Imgur)
	v.Watch(v.model.Node,
		stores.NodeValueChanged,
		stores.NodeDescendantChanged,
		stores.NodeFocus,
		IconEditorChanged,
	)
	return v
}

func (v *IconEditorView) Receive(notif flux.NotifPayload) {
	defer close(notif.Done)
	v.imgur = v.model.Node.Value.(*Imgur)
	vecty.Rerender(v)
	if notif.Type == stores.NodeFocus {
		v.Focus()
	}
}

func (v *IconEditorView) Focus() {
	v.editor.Focus()
}

type iconEditorNotif string

func (iconEditorNotif) IsNotif() {}

const (
	IconEditorChanged iconEditorNotif = "IconEditorChanged"
)

func (v *IconEditorView) Render() *vecty.HTML {
	v.editor = editors.NewStringEditorView(v.Ctx, v.model.Node.Map["url"], editable.Inline)
	url := ""
	if v.imgur.Url != nil {
		url = v.imgur.Url.Value()
	}
	text := "drop image here to upload"
	if i, ok := v.model.Data[iconEditorUploading]; ok && i.(bool) {
		text = "uploading..."
	}
	return elem.Div(
		prop.Class("container-fluid"),
		elem.Div(
			prop.Class("row"),
			elem.Div(
				prop.Class("col-sm-10"),
				vecty.Style("padding-left", "0"),
				vecty.Style("padding-right", "0"),
				v.editor,
			),
			elem.Div(
				prop.Class("col-sm-2"),
				elem.Image(
					prop.Class("img-responsive"),
					style.MaxHeight("200px"),
					prop.Src(url),
				),
			),
		),
		elem.Div(
			prop.Class("row"),
			elem.Div(
				event.DragEnter(func(e *vecty.Event) {}).PreventDefault().StopPropagation(),
				event.DragOver(func(e *vecty.Event) {}).PreventDefault().StopPropagation(),
				event.Drop(func(e *vecty.Event) {
					if e.Get("dataTransfer").Get("files").Get("length").Int() != 1 {
						return
					}
					file := e.Get("dataTransfer").Get("files").Index(0)
					fileType := file.Get("type").String()
					fr := NewFileReader(file)
					go func() {

						v.App.Dispatch(&actions.EditorData{
							Func: func(payload *flux.Payload) {
								v.model.Data[iconEditorUploading] = true
								payload.Notify(v.model.Node, IconEditorChanged)
							},
						})

						req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", fr)
						if err != nil {
							v.App.Fail <- kerr.Wrap("MBXNQMLLTO", err)
						}
						req.Header.Set("Content-Type", fileType)
						req.Header.Set("Authorization", "Client-ID d636cf32baf8829")
						client := &http.Client{}
						resp, err := client.Do(req)
						if err != nil {
							v.App.Fail <- kerr.Wrap("LHCTQKFNOI", err)
						}
						defer resp.Body.Close()

						b, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							v.App.Fail <- kerr.Wrap("IGOTVBWALY", err)
						}
						var body struct {
							Data struct {
								Link string `json:"link"`
							} `json:"data"`
						}
						err = json.Unmarshal(b, &body)
						if err != nil {
							v.App.Fail <- kerr.Wrap("UQIYLQMLFE", err)
						}

						v.App.Dispatch(&actions.EditorData{
							Func: func(payload *flux.Payload) {
								v.model.Data[iconEditorUploading] = false
								payload.Notify(v.model.Node, IconEditorChanged)
							},
						})

						v.App.Dispatch(&actions.Modify{
							Undoer:    &actions.Undoer{},
							Editor:    v.url,
							Before:    v.url.Node.NativeValue(),
							After:     body.Data.Link,
							Immediate: true,
						})

					}()

				}).PreventDefault().StopPropagation(),

				style.Height("100px"),
				vecty.Style("background-color", "#eeeeee"),
				vecty.Text(text),
			),
		),
	)
}
