package assetmap

import "github.com/defrankland/hasherator"

var Assets hasherator.AssetsDir = hasherator.AssetsDir{}

func init() {
	err := Assets.Run("./static/", "./static/hashed/", []string{})
	if err != nil {
		panic(err)
	}
}
