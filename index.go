// Code generated for package switcher by go-bindata DO NOT EDIT. (@generated)
// sources:
// index.html
package switcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _indexHtml = []byte(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
    <link rel="shortcut icon" type="image/x-icon" href="/favicon.ico">
    <title>Mikrotik Provider Switcher</title>
    <link rel="stylesheet" href="/static/bootstrap.min.css">
    <script src="/static/vue.min.js"></script>
    <script src="/static/vue-spinner.min.js"></script>
    <script src="/static/axios.min.js"></script>
    <style scoped>
        .container-fluid {
            margin-right: auto;
            margin-left: auto;
            max-width: 350px; /* or 950px */
        }

        .dune-logo {
            max-height: 34px;
            margin-top: 6px;
        }

        .spacer {
            height: 100px;
        }

        .pulse {
            padding-top: 16px;
        }
    </style>

</head>

<div id="app" class="container-fluid">
    <!-- Switch Providers on Mikrotik -->
    <div id="mikrotik" class="container mt-3 mt-sm-5">
        <div class="container text-center">
            <h2 v-show='!switching'>Провайдеры</h2>
            <pulse-loader class="pulse" :loading="switching" :color="spinnerColor" :size="spinnerSize"></pulse-loader>
        </div>

        <button type="button" class="btn btn-lg btn-block"
                v-bind:class="[ isSumtel ? 'btn-success': 'btn-light']"
                v-on:click="clickSumtel">
            <img src="/static/sumtel_logo.png" alt="Sumtel" about="Sumtel">
        </button>

        <button type="button" class="btn btn-lg btn-block"
                v-bind:class="[ isEtelecom ? 'btn-success': 'btn-light']"
                v-on:click="clickEtelecom">
            <img src="/static/etelecom_logo.png" alt="Etelecom" about="Sumtel">
        </button>

        <div class="alert alert-info mt-2 text-center d-none d-sm-block" role="alert">
        {{ version }}
        </div>
    </div>
    <!-- Dune HD -->
    <div id="dunes" class="container m-2 mt-sm-3">
        <div class="container text-center">
            <h2 v-show='!dune_request'>Dunes</h2>
            <pulse-loader class="pulse"
                          :loading="dune_request"
                          :color="spinnerColor"
                          :size="spinnerSize">
            </pulse-loader>
        </div>
        <div class="row justify-content-center mb-2" v-for="(item, index) in dunes">
            <img class="col-4 dune-logo" src="/static/dune_logo.png"
                 v-bind:alt="'Dune ' + item"
                 v-bind:about="'Dune ' + item">
            <button type="button" class="col-6 btn btn-lg btn-block text-center align-middle"
                    v-bind:id="'dune-' + index"
                    v-bind:disabled="dunes_status[index] === 'offline'"
                    v-bind:class="[dunes_button_class[index]]"
                    v-on:click="clickDune(index)">
                <span>{{ item }}</span>
            </button>
        </div>
    </div>
</div>

<script>
    var PulseLoader = VueSpinner.PulseLoader;

    var vm = new Vue({
        el: '#app',

        components: {
            'PulseLoader': PulseLoader,
        },

        data: {
            version: '',
            provider: '',
            switching: false,
            dune_request: false,
            not_switching: true,
            spinnerColor: '#28a745',
            spinnerSize: '20px',
            dunes: [],
            dunes_status: [],
            dunes_button_class: []
        },

        computed: {
            isSumtel: function () {
                return this.provider === "1";
            },
            isEtelecom: function () {
                return this.provider === "2";
            }
        },

        created: function () {
            console.log('The application has started');
            var self = this;
            this.getVersion().then(function () {
                self.getProvider().then(function () {
                    console.log('Version ' + self.version);
                    console.log('provider ' + self.provider);
                });
            });
            this.getDunes().then(function () {
                console.log('Dunes ' + self.dunes);
                self.dunes_status.length = self.dunes.length;
                self.dunes_status.fill("offline");
                self.dunes_button_class.length = self.dunes.length;
                self.dunes_button_class.fill("btn-danger");

                self.dunes.forEach(function(item, i) {
                    self.getDuneStatus(i)
                });
            })
        },

        methods: {
            //
            // Mikrotik
            //
            getVersion: function () {
                var self = this;
                return axios.get('/api/v1/mikrotik')
                        .then(function (response) {
                            self.version = response.data.version;
                        });
            },
            getProvider: function () {
                var vm = this;
                return axios.get('/api/v1/provider').then(function (response) {
                    vm.provider = response.data.provider;
                });
            },
            switchProvider: function () {
                this.switching = true;
                var vm = this;
                return axios.post('/api/v1/switch').then(function (response) {
                    vm.provider = response.data.provider;
                    vm.switching = false;
                });
            },
            clickEtelecom: function () {
                if (!this.switching && this.isSumtel) this.switchProvider();
            },
            clickSumtel: function () {
                if (!this.switching && this.isEtelecom) this.switchProvider();
            },
            //
            // Dunes
            //
            getDunes: function () {
                var vm = this;
                vm.dune_request = true;
                return axios.get('/api/v1/dune/names').then(function (response) {
                    vm.dunes = response.data.names;
                    vm.dune_request = false;
                });
            },
            getDuneStatus: function (index) {
                var vm = this;
                vm.dune_request = true;
                return axios.get('/api/v1/dune/' + index +'/status').then(function (response) {
                    //console.log('response: ' + response.data.status);
                    vm.dunes_status[index] = response.data.status;
                    console.log('getDuneStatus: ' + index + ' - ' + response.data.status);
                    vm.dune_request = false;
                    vm.refreshDune(index);
                });
            },
            refreshDune: function (index) {
                var res = 'btn-danger';
                if (this.dunes_status[index] === 'offline') {
                    res = 'btn-danger'
                } else {
                    if (this.dunes_status[index] === 'standby') {
                        res = 'btn-secondary'
                    } else {
                        res = 'btn-success'
                    }
                }
                this.dunes_button_class[index] = res
                console.log('refreshDune: ' + index + ' - ' + res);
            },
            clickDune: function (index) {
                var vm = this;
                console.log('Dune button clicked index: ' + index + ' status: ' + vm.dunes_status[index]);
                if (vm.dunes_status[index] === 'standby') {
                    return this.duneOn(index)
                } else {
                    return this.duneOff(index)
                }
            },
            duneOn: function (index) {
                var vm = this;
                console.log('duneOn: ' + index);
                // additional check - may by changed from another device
                vm.getDuneStatus(index).then(function () {
                    if (vm.dunes_status[index] === 'standby') {
                        vm.dune_request = true;
                        axios.get('/api/v1/dune/' + index + '/on').then(function (response) {
                            setTimeout(function () {
                                vm.getDuneStatus(index)
                            }, 5000);
                        })
                    }
                })
            },
            duneOff: function (index) {
                var vm = this;
                if (vm.dunes_status[index] !== 'standby')
                vm.dune_request = true;
                console.log('duneOff: ' + index);
                // additional check - may by changed from another device
                vm.getDuneStatus(index).then(function () {
                    if (vm.dunes_status[index] !== 'standby') {
                        vm.dune_request = true;
                        axios.get('/api/v1/dune/' + index + '/off').then(function (response) {
                            setTimeout(function () {
                                vm.getDuneStatus(index)
                            }, 10000);
                        })
                    }
                })
            }
        }
    });
</script>
</body>
</html>`)

func indexHtmlBytes() ([]byte, error) {
	return _indexHtml, nil
}

func indexHtml() (*asset, error) {
	bytes, err := indexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.html", size: 9270, mode: os.FileMode(420), modTime: time.Unix(1539015318, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"index.html": indexHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"index.html": &bintree{indexHtml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
