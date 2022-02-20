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
	Helptext: cmds.HelpText{
		Tagline: "multiconvert",
		ShortDescription: `
btfs convert 16U...          -> 0x...           (peer id -> bttc addr)
btfs convert CAISI...        -> private key hex (pk base64 -> pk hex)
btfs convert T...            -> 0x...           (tron addr -> bttc addr)
btfs convert 0x...           -> T...            (bttc addr -> tron addr)
btfs convert private key hex -> CAISI...        (pk hex -> pk base64)
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("input", true, false, ""),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		input := req.Arguments[0]
		if strings.HasPrefix(input, "16U") { // peerid -> bttc address
			fmt.Println(convertPeerID2BttcAddr(input))
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
			var addrBytes []byte
			b, _ := hex.DecodeString(crypto.AddressPrefix)
			addrBytes = append(addrBytes, b...)
			b, _ = hex.DecodeString(input[2:])
			addrBytes = append(addrBytes, b...)
			r, _ := crypto.Encode58Check(addrBytes)
			fmt.Println(r)
			return nil
		} else if len(input) == 64 { // private key hex -> base64
			b, _ := hex.DecodeString(input)
			base64.StdEncoding.EncodeToString(b)
			sk, _ := ic.UnmarshalSecp256k1PrivateKey(b)
			b, _ = sk.Bytes()
			fmt.Println(base64.StdEncoding.EncodeToString(b))
			return nil
		}

		return errors.New("error input")
	},
}

func convertPeerID2BttcAddr(peerID string) string {
	tmp, _ := peer.Decode(peerID)
	pppk, _ := tmp.ExtractPublicKey()

	pkBytes, _ := ic.RawFull(pppk)
	pk2, _ := eth.UnmarshalPubkey(pkBytes)

	addr := eth.PubkeyToAddress(*pk2)
	return addr.String()
}
