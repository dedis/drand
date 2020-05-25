package chain

import (
	"time"

	"github.com/drand/drand/key"
	proto "github.com/drand/drand/protobuf/drand"
)

func beaconToProto(b *Beacon) *proto.BeaconPacket {
	return &proto.BeaconPacket{
		PreviousSig: b.PreviousSig,
		Round:       b.Round,
		Signature:   b.Signature,
	}
}

func protoToBeacon(p *proto.BeaconPacket) *Beacon {
	return &Beacon{
		Round:       p.GetRound(),
		Signature:   p.GetSignature(),
		PreviousSig: p.GetPreviousSig(),
	}
}

// InfoFromProto returns a Info from the protocol description
func InfoFromProto(p *proto.ChainInfoPacket) (*Info, error) {
	public := key.KeyGroup.Point()
	if err := public.UnmarshalBinary(p.PublicKey); err != nil {
		return nil, err
	}

	return &Info{
		PublicKey:   public,
		GenesisTime: p.GenesisTime,
		Period:      time.Duration(p.Period) * time.Second,
	}, nil
}
