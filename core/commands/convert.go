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
		switch {
		case strings.HasPrefix(input, "16U"):
			id, _ := peer.Decode(input)
			pk, _ := id.ExtractPublicKey()
			pkBytes, _ := ic.RawFull(pk)
			pk2, _ := eth.UnmarshalPubkey(pkBytes)
			fmt.Println(eth.PubkeyToAddress(*pk2))
			return nil
		case strings.HasPrefix(input, "CAISI"):
			pkb, _ := base64.StdEncoding.DecodeString(input)
			fmt.Println(hex.EncodeToString(pkb[4:]))
			return nil
		case strings.HasPrefix(input, "T"):
			b, _ := crypto.Decode58Check(input)
			fmt.Println(common.BytesToAddress(b))
			return nil
		case strings.HasPrefix(input, "0x"):
			var addrBytes []byte
			b, _ := hex.DecodeString(crypto.AddressPrefix)
			addrBytes = append(addrBytes, b...)
			b, _ = hex.DecodeString(input[2:])
			addrBytes = append(addrBytes, b...)
			r, _ := crypto.Encode58Check(addrBytes)
			fmt.Println(r)
			return nil
		case len(input) == 64:
			b, _ := hex.DecodeString(input)
			base64.StdEncoding.EncodeToString(b)
			sk, _ := ic.UnmarshalSecp256k1PrivateKey(b)
			b, _ = sk.Bytes()
			fmt.Println(base64.StdEncoding.EncodeToString(b))
			return nil
		default:
			return errors.New("error input")
		}
	},
}
