package processing

import (
	"bytes"

	"github.com/trimmer-io/go-xmp/xmp"

	"github.com/imgproxy/imgproxy/v3/imagedata"
	"github.com/imgproxy/imgproxy/v3/imagemeta/iptc"
	"github.com/imgproxy/imgproxy/v3/options"
	"github.com/imgproxy/imgproxy/v3/vips"
)

func stripIPTC(img *vips.Image) []byte {
	iptcData, err := img.GetBlob("iptc-data")
	if err != nil || len(iptcData) == 0 {
		return nil
	}

	iptcMap := make(iptc.IptcMap)
	err = iptc.ParsePS3(iptcData, iptcMap)
	if err != nil {
		return nil
	}

	for key := range iptcMap {
		if key.RecordID == 2 && key.TagID != 80 && key.TagID != 110 && key.TagID != 116 {
			delete(iptcMap, key)
		}
	}

	if len(iptcMap) == 0 {
		return nil
	}

	return iptcMap.Dump()
}

func stripXMP(img *vips.Image) []byte {
	xmpData, err := img.GetBlob("xmp-data")
	if err != nil || len(xmpData) == 0 {
		return nil
	}

	xmpDoc, err := xmp.Read(bytes.NewReader(xmpData))
	if err != nil {
		return nil
	}

	namespaces := xmpDoc.Namespaces()
	filteredNs := namespaces[:0]

	for _, ns := range namespaces {
		if ns.Name == "dc" || ns.Name == "xmpRights" || ns.Name == "cc" {
			filteredNs = append(filteredNs, ns)
		}
	}
	xmpDoc.FilterNamespaces(filteredNs)

	nodes := xmpDoc.Nodes()
	for _, n := range nodes {
		if n.Name() == "dc" {
			filteredNodes := n.Nodes[:0]
			for _, nn := range n.Nodes {
				name := nn.Name()
				if name == "rights" || name == "contributor" || name == "creator" || name == "publisher" {
					filteredNodes = append(filteredNodes, nn)
				}
			}
			n.Nodes = filteredNodes

			filteredAttrs := n.Attr[:0]
			for _, a := range n.Attr {
				name := a.Name.Local
				if name == "dc:rights" || name == "dc:contributor" || name == "dc:creator" || name == "dc:publisher" {
					filteredAttrs = append(filteredAttrs, a)
				}
			}
			n.Attr = filteredAttrs
		}
	}

	if len(xmpDoc.Nodes()) == 0 {
		return nil
	}

	xmpData, err = xmp.Marshal(xmpDoc)
	if err != nil {
		return nil
	}

	return xmpData
}

func stripMetadata(pctx *pipelineContext, img *vips.Image, po *options.ProcessingOptions, imgdata *imagedata.ImageData) error {
	if !po.StripMetadata {
		return nil
	}

	var iptcData, xmpData []byte

	if po.KeepCopyright {
		iptcData = stripIPTC(img)
		xmpData = stripXMP(img)
	}

	if err := img.Strip(po.KeepCopyright); err != nil {
		return err
	}

	if po.KeepCopyright {
		if len(iptcData) > 0 {
			img.SetBlob("iptc-data", iptcData)
		}

		if len(xmpData) > 0 {
			img.SetBlob("xmp-data", xmpData)
		}
	}

	return nil
}
