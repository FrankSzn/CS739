package nfs

import (
	"errors"
	"fmt"
	pb "github.com/Ricky54326/CS739/hw2/protos"
	"os"
	"syscall"
)

func WriteNFS(in *pb.WriteArgs) (*pb.WriteReturn, error) {
	// we have same input args as Read but also have Data which is an array of bytes
	// see read.go for info on some decision made here

	// get path for the file
	filepath, err := InumToPath(int(in.Fh.Inode))
	if err != nil {
		return &pb.WriteReturn{Attr: &pb.Attribute{}}, errors.New("inode not found")
	}

	// ensure genums match
	fh_genum := in.Fh.Genum
	fs_genum, err := PathToGen(filepath)
	if err != nil {
		fmt.Println("Genum of file not found, error.")
		os.Exit(-1)
	}

	if fh_genum != fs_genum {
		return &pb.WriteReturn{Attr: &pb.Attribute{}}, errors.New("genum mismatch")
	}

	// get file object for that file (not an fd)
	f, err := os.OpenFile(filepath, os.O_WRONLY, 0)
	if err != nil {
		return &pb.WriteReturn{Attr: &pb.Attribute{}}, err
	}

	// write the data into it starting at in.Offset
	data_to_write := in.Data[0:in.Count]
	n_bytes_written, err := f.WriteAt(data_to_write, in.Offset)
	n_bytes_written = n_bytes_written // to supress compiler warning

	// get attributes after writing it
	var f_info syscall.Stat_t
	err = syscall.Stat(filepath, &f_info)
	if err != nil {
		fmt.Println("Stat failed, FATAL error")
		os.Exit(-1)
	}

	return &pb.WriteReturn{Attr: StatToAttr(&f_info)}, err
}
