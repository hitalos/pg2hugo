package filesystem

import (
	"context"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/hitalos/pg2hugo/models"
)

type dir struct {
	content *models.Content
	entry   fuse.Dirent
}

func (d dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = d.entry.Inode
	a.Size = 4096
	a.Mode = os.ModeDir | 0o555
	a.Mtime = d.content.LastMod
	a.Ctime = d.content.PublishDate
	return nil
}

func (d dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	folderInode := d.entry.Inode + 1
	files := []fuse.Dirent{
		{Name: ".", Inode: d.entry.Inode, Type: fuse.DT_Dir},
		{Name: "..", Inode: 1, Type: fuse.DT_Dir},
		{Name: "index.md", Inode: folderInode, Type: fuse.DT_File},
	}
	for _, r := range d.content.Resources {
		folderInode++
		files = append(files, fuse.Dirent{Name: r.Src, Inode: folderInode, Type: fuse.DT_File})
	}
	return files, nil
}

func (d dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "index.md" {
		return &index{d.entry.Inode + 1, d.content}, nil
	}
	for i, r := range d.content.Resources {
		if r.Src == name {
			return &file{d.entry.Inode + uint64(i) + 1, r}, nil
		}
	}
	return nil, syscall.ENOENT
}
