package filesystem

import (
	"context"
	"log"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"

	"github.com/hitalos/pg2hugo/models"
)

type file struct {
	inode    uint64
	resource *models.Resource
}

func (f file) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = f.inode
	a.Size = f.resource.Size()
	a.Mode = 0444
	a.Mtime = f.resource.LastMod
	a.Ctime = f.resource.LastMod
	return nil
}

func (f file) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	if len(f.resource.Bs) == 0 {
		if err := f.resource.Load(); err != nil {
			log.Println(err)
			return err
		}
	}
	fuseutil.HandleRead(req, resp, f.resource.Bs)
	return nil
}

func (f file) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	return f, nil
}
