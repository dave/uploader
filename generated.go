// info:{"Path":"github.com/dave/uploader","Hash":15281916385306789409}
package uploader

import (
	context "context"
	fmt "fmt"
	reflect "reflect"

	jsonctx "frizz.io/context/jsonctx"
	system "frizz.io/system"
)

// notest

// Automatically created basic rule for imgur
type ImgurRule struct {
	*system.Object
	*system.Rule
}

func (v *ImgurRule) Unpack(ctx context.Context, in system.Packed, iface bool) error {
	if in == nil || in.Type() == system.J_NULL {
		return nil
	}
	if v.Object == nil {
		v.Object = new(system.Object)
	}
	if err := v.Object.Unpack(ctx, in, false); err != nil {
		return err
	}
	if err := v.Object.InitializeType("github.com/dave/uploader", "@imgur"); err != nil {
		return err
	}
	if v.Rule == nil {
		v.Rule = new(system.Rule)
	}
	if err := v.Rule.Unpack(ctx, in, false); err != nil {
		return err
	}
	return nil
}
func (v *ImgurRule) Repack(ctx context.Context) (data interface{}, typePackage string, typeName string, jsonType system.JsonType, err error) {
	if v == nil {
		return nil, "github.com/dave/uploader", "@imgur", system.J_NULL, nil
	}
	m := map[string]interface{}{}
	if v.Object != nil {
		ob, _, _, _, err := v.Object.Repack(ctx)
		if err != nil {
			return nil, "", "", "", err
		}
		for key, val := range ob.(map[string]interface{}) {
			m[key] = val
		}
	}
	if v.Rule != nil {
		ob, _, _, _, err := v.Rule.Repack(ctx)
		if err != nil {
			return nil, "", "", "", err
		}
		for key, val := range ob.(map[string]interface{}) {
			m[key] = val
		}
	}
	return m, "github.com/dave/uploader", "@imgur", system.J_OBJECT, nil
}

type Imgur struct {
	*system.Object
	Url *system.String `json:"url"`
}
type ImgurInterface interface {
	GetImgur(ctx context.Context) *Imgur
}

func (o *Imgur) GetImgur(ctx context.Context) *Imgur {
	return o
}
func UnpackImgurInterface(ctx context.Context, in system.Packed) (ImgurInterface, error) {
	switch in.Type() {
	case system.J_MAP:
		i, err := system.UnpackUnknownType(ctx, in, true, "github.com/dave/uploader", "imgur")
		if err != nil {
			return nil, err
		}
		ob, ok := i.(ImgurInterface)
		if !ok {
			return nil, fmt.Errorf("%T does not implement ImgurInterface", i)
		}
		return ob, nil
	default:
		return nil, fmt.Errorf("Unsupported json type %s when unpacking into ImgurInterface.", in.Type())
	}
}
func (v *Imgur) Unpack(ctx context.Context, in system.Packed, iface bool) error {
	if in == nil || in.Type() == system.J_NULL {
		return nil
	}
	if v.Object == nil {
		v.Object = new(system.Object)
	}
	if err := v.Object.Unpack(ctx, in, false); err != nil {
		return err
	}
	if err := v.Object.InitializeType("github.com/dave/uploader", "imgur"); err != nil {
		return err
	}
	if field, ok := in.Map()["url"]; ok && field.Type() != system.J_NULL {
		ob0 := new(system.String)
		if err := ob0.Unpack(ctx, field, false); err != nil {
			return err
		}
		v.Url = ob0
	}
	return nil
}
func (v *Imgur) Repack(ctx context.Context) (data interface{}, typePackage string, typeName string, jsonType system.JsonType, err error) {
	if v == nil {
		return nil, "github.com/dave/uploader", "imgur", system.J_NULL, nil
	}
	m := map[string]interface{}{}
	if v.Object != nil {
		ob, _, _, _, err := v.Object.Repack(ctx)
		if err != nil {
			return nil, "", "", "", err
		}
		for key, val := range ob.(map[string]interface{}) {
			m[key] = val
		}
	}
	if v.Url != nil {
		ob0, _, _, _, err := v.Url.Repack(ctx)
		if err != nil {
			return nil, "", "", "", err
		}
		m["url"] = ob0
	}
	return m, "github.com/dave/uploader", "imgur", system.J_OBJECT, nil
}
func init() {
	pkg := jsonctx.InitPackage("github.com/dave/uploader")
	pkg.SetHash(uint64(0xd41445fc416ec621))
	pkg.Init("imgur", func() interface{} {
		return new(Imgur)
	}, nil, func() interface{} {
		return new(ImgurRule)
	}, func() reflect.Type {
		return reflect.TypeOf((*ImgurInterface)(nil)).Elem()
	})
}
