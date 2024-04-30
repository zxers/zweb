package zweb

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/groupcache/lru"
)

type FileUpload struct {
	Key string
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUpload) Handle() HandleFunc {
	return func(ctx *Context) {
		src, fileHeader, err := ctx.Req.FormFile(f.Key)
		if err != nil {
			ctx.RespStatusCode = 400
			ctx.RespData = []byte("未找到数据")
			log.Fatalln(err)
			return
		}
		dst, err := os.OpenFile(f.DstPathFunc(fileHeader), os.O_CREATE, 0o666)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			log.Fatalln(err)
			return
		}
		_, err = io.CopyBuffer(dst, src, nil)
		if err != nil {
			ctx.RespStatusCode = 400
			ctx.RespData = []byte("上传失败")
			log.Fatalln(err)
			return
		}
		ctx.RespStatusCode = 200
		ctx.RespData = []byte("上传成功")
	}
}

type FileDownload struct {
	Dir string
}

func (f *FileDownload) Handle() HandleFunc {
	return func(ctx *Context) {
		req, _ := ctx.QueryValue("file").String()
		path := filepath.Join(f.Dir, filepath.Clean(req))
		log.Println(path)
		fn := filepath.Base(path)
		log.Println(fn)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

type StaticResourceHandler struct {
	dir string
	cache *lru.Cache
	extensionContentTypeMap map[string]string
	maxFileSize int
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

func NewStaticResourceHandler(dir string, opts...StaticResourceHandlerOptions) *StaticResourceHandler {
	s := &StaticResourceHandler{
		dir: dir,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
		maxFileSize: 10 * 1024 * 1024,
		cache: lru.New(1000),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type StaticResourceHandlerOptions func(h *StaticResourceHandler)

func WithMaxFileSize (maxFileSize int)  StaticResourceHandlerOptions {
	return func(h *StaticResourceHandler) {
		h.maxFileSize = maxFileSize
	}
}

func (h *StaticResourceHandler) Handle(ctx *Context) {
	// 从请求路径中获取文件名
	req, _ := ctx.PathValue("file").String()
	// 从缓存中获取
	if item, ok := h.readFileFromData(req); ok {
		log.Printf("从缓存中读取数据...")
		h.writeItemAsResponse(item, ctx.Resp)
		return
	}
	path := filepath.Join(h.dir, req)
	f, err := os.Open(path)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 获取文件类型
	ext := getFileExt(f.Name())
	t, ok := h.extensionContentTypeMap[ext]
	if !ok {
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	item := &fileCacheItem{
		fileSize:    len(data),
		data:        data,
		contentType: t,
		fileName:    req,
	}

	h.cacheFile(item)
	h.writeItemAsResponse(item, ctx.Resp)
}

func (h *StaticResourceHandler) writeItemAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", item.fileSize))
	_, _ = writer.Write(item.data)
}

func (h *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if h.cache != nil && item.fileSize < h.maxFileSize {
		h.cache.Add(item.fileName, item.data)
	}
}

func (h *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if h.cache != nil {
		if val, ok := h.cache.Get(fileName); !ok {
			return val.(*fileCacheItem), ok
		}
	}
	return nil, false
}

func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name)-1 {
		return ""
	}
	return name[index+1:]
}