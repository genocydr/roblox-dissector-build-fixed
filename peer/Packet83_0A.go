package peer

import (
	"errors"
	"fmt"

	"github.com/gskartwii/roblox-dissector/datamodel"
)

// Packet83_0A describes a ID_PROP_ACK packet.
type Packet83_0A struct {
	// Instance that had the property change
	Instance *datamodel.Instance
	Schema   *StaticPropertySchema
	Versions []uint32
}

func (thisBitstream *extendedReader) DecodePacket83_0A(reader PacketReader, layers *PacketLayers) (Packet83Subpacket, error) {
	var err error
	layer := &Packet83_0A{}

	referent, err := thisBitstream.readObject(reader.Caches())
	if err != nil {
		return layer, err
	}
	if referent.IsNull {
		return layer, errors.New("self is null in repl prop ack")
	}
	layer.Instance, err = reader.Context().InstancesByReferent.TryGetInstance(referent)
	if err != nil {
		return layer, err
	}

	context := reader.Context()
	propertyIDx, err := thisBitstream.readUint16BE()
	if err != nil {
		return layer, err
	}

	if int(propertyIDx) > int(len(context.StaticSchema.Properties)) {
		return layer, fmt.Errorf("prop idx %d is higher than %d", propertyIDx, len(context.StaticSchema.Properties))
	}
	layer.Schema = context.StaticSchema.Properties[propertyIDx]

	countVersions, err := thisBitstream.readUint8()
	if err != nil {
		return layer, err
	}
	layer.Versions = make([]uint32, countVersions)
	for i := 0; i < int(countVersions); i++ {
		layer.Versions[i], err = thisBitstream.readUintUTF8()
		if err != nil {
			return layer, err
		}
	}

	return layer, err
}

func (layer *Packet83_0A) Serialize(writer PacketWriter, stream *extendedWriter) error {
	err := stream.writeObject(layer.Instance, writer.Caches())
	if err != nil {
		return err
	}

	context := writer.Context()

	err = stream.writeUint16BE(layer.Schema.NetworkID)
	if err != nil {
		return err
	}
	err = stream.WriteByte(uint8(len(layer.Versions)))
	if err != nil {
		return err
	}
	for i := 0; i < len(layer.Versions); i++ {
		err = stream.writeUintUTF8(layer.Versions[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (Packet83_0A) Type() uint8 {
	return 0xA
}
func (Packet83_0A) TypeString() string {
	return "ID_REPLIC_CFRAME_ACK"
}

func (layer *Packet83_0A) String() string {
	return fmt.Sprintf("ID_REPLIC_CFRAME_ACK: %s[%s]", layer.Instance.GetFullName(), layer.PropertyName)
}
