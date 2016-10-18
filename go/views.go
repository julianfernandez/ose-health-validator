package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/views/index.html": {
		local:   "views/index.html",
		size:    3671,
		modtime: 1476777197,
		compressed: `
H4sIAAAJbogA/+RX31PbOBB+pjP8DxrfI+f4kgAHbZwZCjO0kKO9BtqjLzcbW7GVyJYjyQkhk//9Vv5F
4gRKe9zdw714VtK3n3Y/rVbJYkF8OmQxJRaLfXpnkeVy91Un1BHv7r4yFgU/tzTTnHb7VNKY6Tk552IA
nHwE6JOPUoyop8lpSL0xlR0nx6JTRDWQUOvEppOUTV3rVMSaxtq+nifUIl4+ci1N77Rjdn1DvBCkotpN
9dA+sohT0cQQUdeaMjpLhNQrzjPm69D16ZR51M4GPxOGQTLgtvKAU7dZ8HAWj3dfkVDSoWuZsNRrx4ng
zvPjxkAIrbSExAw8ETnZxAy0FzrtRrtx6CgOmjoVrhExBCplIaOk3LWUnnOqQkq1mWIYXSBRKlwIoXVw
aCepvG/PLk+O0t5e+0LOT9LbcX/8/muvefrubf/+onUzm96dfIhuo99cgj4HzZYd0ON5r381+XB0djw5
ut8b99KbmKd/3F1cf770e6ez+9u3/WP5tU3PvvwqLi73ot7ok/od9o+PD66T8/dXp/vnX27Zfr8lvDOl
0jhsTtqB65oIPSmUEpIFLHYtiEU8j0SqVqTaltbztBvi2dgwo0pE1NlvHDaaDiq1Nl3J1+04ZrNHNy22
VBo08zKaRAo/HUDdHQmUJ1miicbqKopqBFPIZy2ipFfxjJQzmqRUzu1Wo4XHa6IZZWw5uvu9bOt18beo
IpQn1lt5nIcLORD+nHgclHItY2di4Al1fDYt52eUc2I+toqyRdIJm8+8xAg0bA7SZRvWqM31A+wc0h7y
lPkF+wpAilk+WXPjduTbzVaxRjopL9dimOItn9oJ41xl1ihVmg0Z9Us0wRop4eBpNqWoDxAfNNhaBAE3
8sKgLJqfUC6uQ8Tg0Yg46L7LxuQzcIY+AtMsFjoOZJW0ss/TxGB1ScV68m3GjpPmPZXUFEFWu+hlZZZl
DAbFfLxueRarHglg0x6CT7HRkFKJwu/b57SBwsMiASQrq0+f2gbCxMNJ9rXxSYGU6zXsJtg2pcziYB2G
Bdpeh2WPSQ1UKf9JpJqqSvN1Jidsr4dQ1PJTQZmbVNtssZAQB5Q0zNtItjsXArXrceIDy4aETkijqDxL
jPNHltSZSioOA1Qy+9oq9TyqFL51UWKt+5CsRbqWz1TCYf6aDLjwxm8sUosApcLMSvIhYNHYnrnjxhjO
zJcHWZNBWDdHm8rP6nyxaNxIvlwSDTLAR9n6c8AhHltdXLjC93i5zKs8c6rlZMQ2wezs7CwWNPbrSRcr
dX3G4jF9tgnkm5ORL62PZhFVT+mzkv5jqdcy31aHOzumQDJAocgG5JHx+rDie7AqY6OVwPY28r9sHjTh
zIN/r32QH+wfxrsMlrQbv/yzDeRFOkh1Q4iNSVTRf8eNqboD1xvZk2em/1+3hx9Mnqz1gi2Ft9lQX6ph
1MAPZmWVRscxxZ7/FnWKf4xFXH8FAAD//4Y4X05XDgAA
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/views": {
		isDir: true,
		local: "/views",
	},
}
