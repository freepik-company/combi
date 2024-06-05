package nginx

// ----------------------------------------------------------------
// Merge NGINX data structure
// ----------------------------------------------------------------
func (e *NginxT) GetConfigStruct() (config interface{}) {
	return &e.ConfigStruct
}

func (e *NginxT) MergeConfigs(source interface{}) {
	mergeNginxBlockContent(&e.ConfigStruct, source.(*BlockContentT))
}

func mergeNginxBlockContent(destination, source *BlockContentT) {
	for _, srcVal := range source.Directives {
		found := false
		for dstIn, dstVal := range destination.Directives {
			if dstVal.Name == srcVal.Name && dstVal.Param == srcVal.Param {
				destination.Directives[dstIn] = srcVal
				found = true
			}
		}
		if !found {
			destination.Directives = append(destination.Directives, srcVal)
		}
	}

	for _, srcVal := range source.Blocks {
		found := false
		for dstIn, dstVal := range destination.Blocks {
			if dstVal.Name == srcVal.Name && dstVal.Params == srcVal.Params {
				mergeNginxBlockContent(&destination.Blocks[dstIn].BlockContent, &srcVal.BlockContent)
				found = true
			}
		}
		if !found {
			destination.Blocks = append(destination.Blocks, srcVal)
		}
	}
}
