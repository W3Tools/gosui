package secp256k1

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

// deterministicSign generates a deterministic ECDSA signature using RFC 6979 for the secp256k1 curve.
func deterministicSign(priv *btcec.PrivateKey, hash []byte) (*big.Int, *big.Int, error) {
	curve := btcec.S256()
	k := deterministicK(priv, hash)
	r, s, err := signWithK(priv, hash, k)
	if err != nil {
		return nil, nil, err
	}

	// Ensure s is in the lower half of the order
	halfOrder := new(big.Int).Rsh(curve.Params().N, 1)
	if s.Cmp(halfOrder) > 0 {
		s.Sub(curve.Params().N, s)
	}

	return r, s, nil
}

// deterministicK generates a deterministic K value according to RFC 6979 for the secp256k1 curve.
func deterministicK(priv *btcec.PrivateKey, hash []byte) *big.Int {
	curve := btcec.S256()
	q := curve.Params().N

	v := make([]byte, 32)
	k := make([]byte, 32)
	for i := range v {
		v[i] = 0x01
	}

	hm := hmac.New(sha256.New, k)
	hm.Write(v)
	hm.Write([]byte{0x00})
	hm.Write(priv.ToECDSA().D.Bytes())
	hm.Write(hash)
	k = hm.Sum(nil)

	hm = hmac.New(sha256.New, k)
	hm.Write(v)
	v = hm.Sum(nil)

	hm = hmac.New(sha256.New, k)
	hm.Write(v)
	hm.Write([]byte{0x01})
	hm.Write(priv.ToECDSA().D.Bytes())
	hm.Write(hash)
	k = hm.Sum(nil)

	hm = hmac.New(sha256.New, k)
	hm.Write(v)
	v = hm.Sum(nil)

	for {
		hm = hmac.New(sha256.New, k)
		hm.Write(v)
		v = hm.Sum(nil)
		kInt := new(big.Int).SetBytes(v)
		if kInt.Sign() > 0 && kInt.Cmp(q) < 0 {
			return kInt
		}
		hm = hmac.New(sha256.New, k)
		hm.Write(v)
		hm.Write([]byte{0x00})
		k = hm.Sum(nil)

		hm = hmac.New(sha256.New, k)
		hm.Write(v)
		v = hm.Sum(nil)
	}
}

// signWithK creates an ECDSA signature using the provided private key, hash, and deterministic K value for the secp256k1 curve.
func signWithK(priv *btcec.PrivateKey, hash []byte, k *big.Int) (*big.Int, *big.Int, error) {
	curve := btcec.S256()
	n := curve.Params().N

	kInv := new(big.Int).ModInverse(k, n)
	r, _ := curve.ScalarBaseMult(k.Bytes())
	r.Mod(r, n)
	if r.Sign() == 0 {
		return nil, nil, fmt.Errorf("calculated R is zero")
	}

	e := new(big.Int).SetBytes(hash)
	s := new(big.Int).Mul(priv.ToECDSA().D, r)
	s.Add(s, e)
	s.Mul(s, kInv)
	s.Mod(s, n)
	if s.Sign() == 0 {
		return nil, nil, fmt.Errorf("calculated S is zero")
	}

	return r, s, nil
}
