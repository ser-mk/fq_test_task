package raw

import "C"
import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"strconv"
	"time"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG,
		Description: "pg_control",
		Groups:      []string{format.PROBE},
		DecodeFn:    decodePg,
	})
}

var stateMap = scalar.UToSymStr{
	0: "DB_STARTUP",
	1: "DB_SHUTDOWNED",
	2: "DB_SHUTDOWNED_IN_RECOVERY",
	3: "DB_SHUTDOWNING",
	4: "DB_IN_CRASH_RECOVERY",
	5: "DB_IN_ARCHIVE_RECOVERY",
	6: "DB_IN_PRODUCTION",
}

var full_page_writesMap = scalar.UToSymStr{
	0: "off",
	1: "on",
}

var timestampMapper = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	ts, ok := s.Actual.(uint64)
	if !ok {
		return s, nil
	}
	s.Sym = time.Unix(int64(ts), 0).UTC().String()
	return s, nil
})

var NextXIDMapper = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	ts, ok := s.Actual.(uint64)
	if !ok {
		return s, nil
	}
	s.Sym = strconv.Itoa(int(ts>>32)) + ":" + strconv.Itoa(int(ts&0xFFFFFFFF))
	return s, nil
})

var halfHexMapper = scalar.Fn(func(s scalar.S) (scalar.S, error) {
	ts, ok := s.Actual.(uint64)
	if !ok {
		return s, nil
	}
	s.Sym = strconv.Itoa(int(ts>>32)) + "/" + strconv.FormatInt(int64(ts&0xFFFFFFFF), 16)
	return s, nil
})

func decodePg(d *decode.D, in any) any {
	d.FieldU64LE("Database system identifier")
	d.FieldU32LE("pg_control version number")
	d.FieldU32LE("Catalog version number")
	d.FieldU64LE("Database cluster state", stateMap)
	d.FieldU64LE("pg_control last modified", timestampMapper)
	d.FieldU64LE("Latest checkpoint location", halfHexMapper)
	d.FieldStruct("copy of last check point record", func(d *decode.D) {
		d.FieldU64LE("Latest checkpoint's REDO location", halfHexMapper)
		d.FieldU32LE("Latest checkpoint's TimeLineID")
		d.FieldU32LE("Latest checkpoint's PrevTimeLineID")
		d.FieldU64LE("Latest checkpoint's full_page_writes", full_page_writesMap)
		d.FieldU64LE("Latest checkpoint's NextXID", NextXIDMapper)
		d.FieldU32LE("Latest checkpoint's NextOID")
		d.FieldU32LE("Latest checkpoint's NextMultiXactId")
		d.FieldU32LE("Latest checkpoint's NextMultiOffset")
		d.FieldU32LE("Latest checkpoint's oldestXID")
		d.FieldU32LE("Latest checkpoint's oldestXID's DB")

		d.FieldU32LE("Latest checkpoint's oldestMultiXid")
		d.FieldU64LE("Latest checkpoint's oldestMulti's DB")
		d.FieldU64LE("Time of latest checkpoint", timestampMapper)
		d.FieldU64LE("Latest checkpoint's oldestCommitTsXid")
		d.FieldU32LE("Latest checkpoint's newestCommitTsXid")
		d.FieldU32LE("Latest checkpoint's oldestActiveXID")
	})
	d.FieldU64LE("Fake LSN counter for unlogged rels", halfHexMapper)
	d.FieldU64LE("Minimum recovery ending location", halfHexMapper)

	//i, err := strconv.ParseInt("1658172234", 10, 64)
	//if err != nil {
	//	panic(err)
	//}
	//tm := time.Unix(i, 0)
	//d.FieldU64LE("Latest checkpoint location")
	//d.FieldU64LE("checkPointCopy")
	//d.FieldU64LE("unloggedLSN")
	//d.FieldU64LE("minRecoveryPoint")
	//d.FieldU32LE("minRecoveryPointTLI")
	return nil
}
