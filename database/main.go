package database

import (
	"github.com/dgraph-io/badger/v4"

	"x.realy.lol/chk"
	"x.realy.lol/log"
	"x.realy.lol/units"
)

type D struct {
	dataDir        string
	BlockCacheSize int
	Logger         *logger
	InitLogLevel   int
	// DB is the badger db
	*badger.DB
	// seq is the monotonic collision free index for raw event storage.
	seq *badger.Sequence
}

func New() (d *D) {
	d = &D{BlockCacheSize: units.Gb}
	return
}

// Path returns the path where the database files are stored.
func (d *D) Path() string { return d.dataDir }

// Init sets up the database with the loaded configuration.
func (d *D) Init(path string) (err error) {
	d.dataDir = path
	log.I.Ln("opening realy event store at", d.dataDir)
	opts := badger.DefaultOptions(d.dataDir)
	opts.BlockCacheSize = int64(d.BlockCacheSize)
	opts.BlockSize = units.Gb
	opts.CompactL0OnClose = true
	opts.LmaxCompaction = true
	d.Logger = NewLogger(d.InitLogLevel, d.dataDir)
	opts.Logger = d.Logger
	if d.DB, err = badger.Open(opts); chk.E(err) {
		return err
	}
	log.T.Ln("getting event store sequence index", d.dataDir)
	if d.seq, err = d.DB.GetSequence([]byte("events"), 1000); chk.E(err) {
		return err
	}
	return nil

}

func (d *D) Close() (err error) { return d.DB.Close() }

// Serial returns the next monotonic conflict free unique serial on the database.
func (d *D) Serial() (ser uint64, err error) {
	if ser, err = d.seq.Next(); chk.E(err) {
	}
	// log.T.ToSliceOfBytes("serial %x", ser)
	return
}
