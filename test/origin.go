package test

import (
	"fmt"
)

// Origin describes the "source" of a [Test].
type Origin interface {
	Path() string
	String() string
	AcceptVisitor(OriginVisitor, ...VisitOption)
}

// OriginVisitor is an interface for dispatching based on the concrete type of
// an [Origin].
type OriginVisitor interface {
	VisitDirectoryOrigin(DirectoryOrigin)
	VisitFileOrigin(FileOrigin)
	VisitFileLineOrigin(FileLineOrigin)
}

// DirectoryOrigin is an [Origin] that refers to a filesystem directory.
type DirectoryOrigin struct {
	DirPath string
}

// Path returns the path to the directory.
func (o DirectoryOrigin) Path() string {
	return string(o.DirPath)
}

func (o DirectoryOrigin) String() string {
	return string(o.DirPath + "/")
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o DirectoryOrigin) AcceptVisitor(v OriginVisitor, options ...VisitOption) {
	cfg := newVisitConfig(options)
	if cfg.TestingT != nil {
		cfg.TestingT.Helper()
	}
	v.VisitDirectoryOrigin(o)
}

// DapperString returns the string used to represent o in
// [github.com/dogmatiq/dapper] output.
func (o DirectoryOrigin) DapperString() string {
	return o.String()
}

// FileOrigin is an [Origin] that refers to a specific file.
type FileOrigin struct {
	FilePath string
}

// Path returns the path to the document file.
func (o FileOrigin) Path() string {
	return string(o.FilePath)
}

func (o FileOrigin) String() string {
	return string(o.FilePath)
}

// DapperString returns the string used to represent o in
// [github.com/dogmatiq/dapper] output.
func (o FileOrigin) DapperString() string {
	return o.String()
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o FileOrigin) AcceptVisitor(v OriginVisitor, options ...VisitOption) {
	cfg := newVisitConfig(options)
	if cfg.TestingT != nil {
		cfg.TestingT.Helper()
	}
	v.VisitFileOrigin(o)
}

// FileLineOrigin is an [Origin] that describes a specific line within a file.
type FileLineOrigin struct {
	FilePath string
	Line     int
}

// Path returns the path to the file.
func (o FileLineOrigin) Path() string {
	return o.FilePath
}

func (o FileLineOrigin) String() string {
	return fmt.Sprintf("%s:%d", o.FilePath, o.Line)
}

// DapperString returns the string used to represent o in
// [github.com/dogmatiq/dapper] output.
func (o FileLineOrigin) DapperString() string {
	return o.String()
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o FileLineOrigin) AcceptVisitor(v OriginVisitor, options ...VisitOption) {
	cfg := newVisitConfig(options)
	if cfg.TestingT != nil {
		cfg.TestingT.Helper()
	}
	v.VisitFileLineOrigin(o)
}
