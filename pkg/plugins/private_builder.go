package plugins

import (
	"fmt"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

func (b *Builder) LoadPrivatePlugins(goPath string, privatePlugins map[string]Descriptor) error {
	for pName, privatePlugin := range privatePlugins {
		manifest, err := ReadManifest(goPath, privatePlugin.ModuleName)
		if err != nil {
			return fmt.Errorf("%s: failed to read private manifest: %w", privatePlugin.ModuleName, err)
		}

		i := interp.New(interp.Options{GoPath: goPath})
		i.Use(stdlib.Symbols)
		i.Use(unsafe.Symbols)

		_, err = i.Eval(fmt.Sprintf(`import "%s"`, manifest.Import))
		if err != nil {
			return fmt.Errorf("%s: failed to import private plugin code %q: %w", privatePlugin.ModuleName, manifest.Import, err)
		}

		b.descriptors[pName] = pluginContext{
			interpreter: i,
			GoPath:      goPath,
			Import:      manifest.Import,
			BasePkg:     manifest.BasePkg,
		}
	}
	return nil
}
