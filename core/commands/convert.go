package commands

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
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
		if strings.HasPrefix(req.Arguments[0], "16U") {
			fmt.Println(ConvertPeerID2BttcAddr(req.Arguments[0]))
			return nil
		} else if strings.HasPrefix(req.Arguments[0], "CAISI") {
			pkb, _ := base64.StdEncoding.DecodeString(req.Arguments[0])
			fmt.Println(hex.EncodeToString(pkb[4:]))
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
