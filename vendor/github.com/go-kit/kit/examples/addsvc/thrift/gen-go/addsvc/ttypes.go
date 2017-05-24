// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package addsvc

import (
	"bytes"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

// Attributes:
//  - Value
//  - Err
type SumReply struct {
	Value int64  `thrift:"value,1" json:"value"`
	Err   string `thrift:"err,2" json:"err"`
}

func NewSumReply() *SumReply {
	return &SumReply{}
}

func (p *SumReply) GetValue() int64 {
	return p.Value
}

func (p *SumReply) GetErr() string {
	return p.Err
}
func (p *SumReply) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *SumReply) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Value = v
	}
	return nil
}

func (p *SumReply) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Err = v
	}
	return nil
}

func (p *SumReply) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("SumReply"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *SumReply) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("value", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:value: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.Value)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.value (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:value: ", p), err)
	}
	return err
}

func (p *SumReply) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("err", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:err: ", p), err)
	}
	if err := oprot.WriteString(string(p.Err)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.err (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:err: ", p), err)
	}
	return err
}

func (p *SumReply) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SumReply(%+v)", *p)
}

// Attributes:
//  - Value
//  - Err
type ConcatReply struct {
	Value string `thrift:"value,1" json:"value"`
	Err   string `thrift:"err,2" json:"err"`
}

func NewConcatReply() *ConcatReply {
	return &ConcatReply{}
}

func (p *ConcatReply) GetValue() string {
	return p.Value
}

func (p *ConcatReply) GetErr() string {
	return p.Err
}
func (p *ConcatReply) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *ConcatReply) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Value = v
	}
	return nil
}

func (p *ConcatReply) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Err = v
	}
	return nil
}

func (p *ConcatReply) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ConcatReply"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *ConcatReply) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("value", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:value: ", p), err)
	}
	if err := oprot.WriteString(string(p.Value)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.value (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:value: ", p), err)
	}
	return err
}

func (p *ConcatReply) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("err", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:err: ", p), err)
	}
	if err := oprot.WriteString(string(p.Err)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.err (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:err: ", p), err)
	}
	return err
}

func (p *ConcatReply) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ConcatReply(%+v)", *p)
}
