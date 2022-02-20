package commands

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tron-us/go-btfs-common/crypto"
	"strings"

	cmds "github.com/bittorrent/go-btfs-cmds"
	eth "github.com/ethereum/go-ethereum/crypto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

type ConvertOutput struct {
	BttcAddress string
}

var ConvertCmd = &cmds.Command{
	Arguments: []cmds.Argument{
		cmds.StringArg("input", true, false, ""),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		input := req.Arguments[0]
		if strings.HasPrefix(input, "16U") { // peerid -> bttc address
			fmt.Println(ConvertPeerID2BttcAddr(input))
			return nil
		} else if strings.HasPrefix(input, "CAISI") { // private key base64 -> hex
			pkb, _ := base64.StdEncoding.DecodeString(input)
			fmt.Println(hex.EncodeToString(pkb[4:]))
			return nil
		} else if strings.HasPrefix(input, "T") { // tron address -> bttc address
			b, _ := crypto.Decode58Check(input)
			fmt.Println(common.BytesToAddress(b))
			return nil
		} else if strings.HasPrefix(input, "0x") { // bttc address -> tron address
			fmt.Println("TODO")
			return nil
		} else if len(input) == 64 { // private key hex -> base64
			fmt.Println("TODO")
			return nil
		}

		return errors.New("error input")
	},
}

func ConvertPeerID2BttcAddr(peerID string) string {
	tmp, _ := peer.Decode(peerID)
	pppk, _ := tmp.ExtractPublicKey()

	pkBytes, _ := ic.RawFull(pppk)
	pk2, _ := eth.UnmarshalPubkey(pkBytes)

	addr := eth.PubkeyToAddress(*pk2)
	return addr.String()
}
