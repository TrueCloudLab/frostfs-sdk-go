package session_test

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	v2session "github.com/TrueCloudLab/frostfs-api-go/v2/session"
	cidtest "github.com/TrueCloudLab/frostfs-sdk-go/container/id/test"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	frostfsecdsa "github.com/TrueCloudLab/frostfs-sdk-go/crypto/ecdsa"
	oidtest "github.com/TrueCloudLab/frostfs-sdk-go/object/id/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
	sessiontest "github.com/TrueCloudLab/frostfs-sdk-go/session/test"
	"github.com/TrueCloudLab/frostfs-sdk-go/user"
	usertest "github.com/TrueCloudLab/frostfs-sdk-go/user/test"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/util/slice"
	"github.com/stretchr/testify/require"
)

func randSigner() ecdsa.PrivateKey {
	k, err := keys.NewPrivateKey()
	if err != nil {
		panic(fmt.Sprintf("generate private key: %v", err))
	}

	return k.PrivateKey
}

func randPublicKey() frostfscrypto.PublicKey {
	k := randSigner().PublicKey
	return (*frostfsecdsa.PublicKey)(&k)
}

func TestObjectProtocolV2(t *testing.T) {
	var validV2 v2session.Token

	var body v2session.TokenBody
	validV2.SetBody(&body)

	// ID
	id := uuid.New()
	binID, err := id.MarshalBinary()
	require.NoError(t, err)
	restoreID := func() {
		body.SetID(binID)
	}
	restoreID()

	// Owner
	usr := *usertest.ID()
	var usrV2 refs.OwnerID
	usr.WriteToV2(&usrV2)
	restoreUser := func() {
		body.SetOwnerID(&usrV2)
	}
	restoreUser()

	// Lifetime
	var lifetime v2session.TokenLifetime
	lifetime.SetIat(1)
	lifetime.SetNbf(2)
	lifetime.SetExp(3)
	restoreLifetime := func() {
		body.SetLifetime(&lifetime)
	}
	restoreLifetime()

	// Session key
	signer := randSigner()
	authKey := frostfsecdsa.PublicKey(signer.PublicKey)
	binAuthKey := make([]byte, authKey.MaxEncodedSize())
	binAuthKey = binAuthKey[:authKey.Encode(binAuthKey)]
	restoreAuthKey := func() {
		body.SetSessionKey(binAuthKey)
	}
	restoreAuthKey()

	// Context
	cnr := cidtest.ID()
	obj1 := oidtest.ID()
	obj2 := oidtest.ID()
	var cnrV2 refs.ContainerID
	cnr.WriteToV2(&cnrV2)
	var obj1V2 refs.ObjectID
	obj1.WriteToV2(&obj1V2)
	var obj2V2 refs.ObjectID
	obj2.WriteToV2(&obj2V2)
	var cObj v2session.ObjectSessionContext
	restoreCtx := func() {
		cObj.SetTarget(&cnrV2, obj1V2, obj2V2)
		body.SetContext(&cObj)
	}
	restoreCtx()

	// Signature
	var sig refs.Signature
	restoreSig := func() {
		validV2.SetSignature(&sig)
	}
	restoreSig()

	// TODO(@cthulhu-rider): #260 use functionality for message corruption

	for _, testcase := range []struct {
		name      string
		corrupt   []func()
		restore   func()
		assert    func(session.Object)
		breakSign func(*v2session.Token)
	}{
		{
			name: "Signature",
			corrupt: []func(){
				func() {
					validV2.SetSignature(nil)
				},
			},
			restore: restoreSig,
		},
		{
			name: "ID",
			corrupt: []func(){
				func() {
					body.SetID([]byte{1, 2, 3})
				},
				func() {
					id, err := uuid.NewDCEPerson()
					require.NoError(t, err)
					bindID, err := id.MarshalBinary()
					require.NoError(t, err)
					body.SetID(bindID)
				},
			},
			restore: restoreID,
			assert: func(val session.Object) {
				require.Equal(t, id, val.ID())
			},
			breakSign: func(m *v2session.Token) {
				id := m.GetBody().GetID()
				id[len(id)-1]++
			},
		},
		{
			name: "User",
			corrupt: []func(){
				func() {
					var brokenUsrV2 refs.OwnerID
					brokenUsrV2.SetValue(append(usrV2.GetValue(), 1))
					body.SetOwnerID(&brokenUsrV2)
				},
			},
			restore: restoreUser,
			assert: func(val session.Object) {
				require.Equal(t, usr, val.Issuer())
			},
			breakSign: func(m *v2session.Token) {
				id := m.GetBody().GetOwnerID().GetValue()
				copy(id, usertest.ID().WalletBytes())
			},
		},
		{
			name: "Lifetime",
			corrupt: []func(){
				func() {
					body.SetLifetime(nil)
				},
			},
			restore: restoreLifetime,
			assert: func(val session.Object) {
				require.True(t, val.InvalidAt(1))
				require.False(t, val.InvalidAt(2))
				require.True(t, val.InvalidAt(3))
			},
			breakSign: func(m *v2session.Token) {
				lt := m.GetBody().GetLifetime()
				lt.SetIat(lt.GetIat() + 1)
			},
		},
		{
			name: "Auth key",
			corrupt: []func(){
				func() {
					body.SetSessionKey(nil)
				},
				func() {
					body.SetSessionKey([]byte{})
				},
			},
			restore: restoreAuthKey,
			assert: func(val session.Object) {
				require.True(t, val.AssertAuthKey(&authKey))
			},
			breakSign: func(m *v2session.Token) {
				body := m.GetBody()
				key := body.GetSessionKey()
				cp := slice.Copy(key)
				cp[len(cp)-1]++
				body.SetSessionKey(cp)
			},
		},
		{
			name: "Context",
			corrupt: []func(){
				func() {
					body.SetContext(nil)
				},
				func() {
					cObj.SetTarget(nil)
				},
				func() {
					var brokenCnr refs.ContainerID
					brokenCnr.SetValue(append(cnrV2.GetValue(), 1))
					cObj.SetTarget(&brokenCnr)
				},
				func() {
					var brokenObj refs.ObjectID
					brokenObj.SetValue(append(obj1V2.GetValue(), 1))
					cObj.SetTarget(&cnrV2, brokenObj)
				},
			},
			restore: restoreCtx,
			assert: func(val session.Object) {
				require.True(t, val.AssertContainer(cnr))
				require.False(t, val.AssertContainer(cidtest.ID()))
				require.True(t, val.AssertObject(obj1))
				require.True(t, val.AssertObject(obj2))
				require.False(t, val.AssertObject(oidtest.ID()))
			},
			breakSign: func(m *v2session.Token) {
				cnr := m.GetBody().GetContext().(*v2session.ObjectSessionContext).GetContainer().GetValue()
				cnr[len(cnr)-1]++
			},
		},
	} {
		var val session.Object

		for i, corrupt := range testcase.corrupt {
			corrupt()
			require.Error(t, val.ReadFromV2(validV2), testcase.name, fmt.Sprintf("corrupt #%d", i))

			testcase.restore()
			require.NoError(t, val.ReadFromV2(validV2), testcase.name, fmt.Sprintf("corrupt #%d", i))

			if testcase.assert != nil {
				testcase.assert(val)
			}

			if testcase.breakSign != nil {
				require.NoError(t, val.Sign(signer), testcase.name)
				require.True(t, val.VerifySignature(), testcase.name)

				var signedV2 v2session.Token
				val.WriteToV2(&signedV2)

				var restored session.Object
				require.NoError(t, restored.ReadFromV2(signedV2), testcase.name)
				require.True(t, restored.VerifySignature(), testcase.name)

				testcase.breakSign(&signedV2)

				require.NoError(t, restored.ReadFromV2(signedV2), testcase.name)
				require.False(t, restored.VerifySignature(), testcase.name)
			}
		}
	}
}

func TestObject_WriteToV2(t *testing.T) {
	var val session.Object

	assert := func(baseAssert func(v2session.Token)) {
		var m v2session.Token
		val.WriteToV2(&m)
		baseAssert(m)
	}

	// ID
	id := uuid.New()

	binID, err := id.MarshalBinary()
	require.NoError(t, err)

	val.SetID(id)
	assert(func(m v2session.Token) {
		require.Equal(t, binID, m.GetBody().GetID())
	})

	// Owner/Signature
	signer := randSigner()

	require.NoError(t, val.Sign(signer))

	var usr user.ID
	user.IDFromKey(&usr, signer.PublicKey)

	var usrV2 refs.OwnerID
	usr.WriteToV2(&usrV2)

	assert(func(m v2session.Token) {
		require.Equal(t, &usrV2, m.GetBody().GetOwnerID())

		sig := m.GetSignature()
		require.NotZero(t, sig.GetKey())
		require.NotZero(t, sig.GetSign())
	})

	// Lifetime
	const iat, nbf, exp = 1, 2, 3
	val.SetIat(iat)
	val.SetNbf(nbf)
	val.SetExp(exp)

	assert(func(m v2session.Token) {
		lt := m.GetBody().GetLifetime()
		require.EqualValues(t, iat, lt.GetIat())
		require.EqualValues(t, nbf, lt.GetNbf())
		require.EqualValues(t, exp, lt.GetExp())
	})

	// Context
	assert(func(m v2session.Token) {
		cCnr, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Zero(t, cCnr.GetContainer())
		require.Zero(t, cCnr.GetObjects())
	})

	cnr := cidtest.ID()

	var cnrV2 refs.ContainerID
	cnr.WriteToV2(&cnrV2)

	obj1 := oidtest.ID()
	obj2 := oidtest.ID()

	var obj1V2 refs.ObjectID
	obj1.WriteToV2(&obj1V2)
	var obj2V2 refs.ObjectID
	obj2.WriteToV2(&obj2V2)

	val.BindContainer(cnr)
	val.LimitByObjects(obj1, obj2)

	assert(func(m v2session.Token) {
		cCnr, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Equal(t, &cnrV2, cCnr.GetContainer())
		require.Equal(t, []refs.ObjectID{obj1V2, obj2V2}, cCnr.GetObjects())
	})
}

func TestObject_BindContainer(t *testing.T) {
	var val session.Object
	var m v2session.Token
	filled := sessiontest.Object()

	assertDefaults := func() {
		cCnr, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Zero(t, cCnr.GetContainer())
		require.Zero(t, cCnr.GetObjects())
	}

	assertBinary := func(baseAssert func()) {
		val2 := filled

		require.NoError(t, val2.Unmarshal(val.Marshal()))
		baseAssert()
	}

	assertJSON := func(baseAssert func()) {
		val2 := filled

		jd, err := val.MarshalJSON()
		require.NoError(t, err)

		require.NoError(t, val2.UnmarshalJSON(jd))
		baseAssert()
	}

	val.WriteToV2(&m)

	assertDefaults()
	assertBinary(assertDefaults)
	assertJSON(assertDefaults)

	// set value

	cnr := cidtest.ID()

	var cnrV2 refs.ContainerID
	cnr.WriteToV2(&cnrV2)

	val.BindContainer(cnr)

	val.WriteToV2(&m)

	assertCnr := func() {
		cObj, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Equal(t, &cnrV2, cObj.GetContainer())
	}

	assertCnr()
	assertBinary(assertCnr)
	assertJSON(assertCnr)
}

func TestObject_AssertContainer(t *testing.T) {
	var x session.Object

	cnr := cidtest.ID()

	require.False(t, x.AssertContainer(cnr))

	x.BindContainer(cnr)

	require.True(t, x.AssertContainer(cnr))
}

func TestObject_LimitByObjects(t *testing.T) {
	var val session.Object
	var m v2session.Token
	filled := sessiontest.Object()

	assertDefaults := func() {
		cCnr, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Zero(t, cCnr.GetContainer())
		require.Zero(t, cCnr.GetObjects())
	}

	assertBinary := func(baseAssert func()) {
		val2 := filled

		require.NoError(t, val2.Unmarshal(val.Marshal()))
		baseAssert()
	}

	assertJSON := func(baseAssert func()) {
		val2 := filled

		jd, err := val.MarshalJSON()
		require.NoError(t, err)

		require.NoError(t, val2.UnmarshalJSON(jd))
		baseAssert()
	}

	val.WriteToV2(&m)

	assertDefaults()
	assertBinary(assertDefaults)
	assertJSON(assertDefaults)

	// set value

	obj1 := oidtest.ID()
	obj2 := oidtest.ID()

	var obj1V2 refs.ObjectID
	obj1.WriteToV2(&obj1V2)
	var obj2V2 refs.ObjectID
	obj2.WriteToV2(&obj2V2)

	val.LimitByObjects(obj1, obj2)

	val.WriteToV2(&m)

	assertObj := func() {
		cObj, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Equal(t, []refs.ObjectID{obj1V2, obj2V2}, cObj.GetObjects())
	}

	assertObj()
	assertBinary(assertObj)
	assertJSON(assertObj)
}

func TestObject_AssertObject(t *testing.T) {
	var x session.Object

	obj1 := oidtest.ID()
	obj2 := oidtest.ID()
	objOther := oidtest.ID()

	require.True(t, x.AssertObject(obj1))
	require.True(t, x.AssertObject(obj2))
	require.True(t, x.AssertObject(objOther))

	x.LimitByObjects(obj1, obj2)

	require.True(t, x.AssertObject(obj1))
	require.True(t, x.AssertObject(obj2))
	require.False(t, x.AssertObject(objOther))
}

func TestObject_InvalidAt(t *testing.T) {
	var x session.Object

	nbf := rand.Uint64()
	if nbf == math.MaxUint64 {
		nbf--
	}

	iat := nbf
	exp := iat + 1

	x.SetNbf(nbf)
	x.SetIat(iat)
	x.SetExp(exp)

	require.True(t, x.InvalidAt(nbf-1))
	require.True(t, x.InvalidAt(iat-1))
	require.False(t, x.InvalidAt(iat))
	require.True(t, x.InvalidAt(exp))
}

func TestObject_ID(t *testing.T) {
	var x session.Object

	require.Zero(t, x.ID())

	id := uuid.New()

	x.SetID(id)

	require.Equal(t, id, x.ID())
}

func TestObject_AssertAuthKey(t *testing.T) {
	var x session.Object

	key := randPublicKey()

	require.False(t, x.AssertAuthKey(key))

	x.SetAuthKey(key)

	require.True(t, x.AssertAuthKey(key))
}

func TestObject_ForVerb(t *testing.T) {
	var val session.Object
	var m v2session.Token
	filled := sessiontest.Object()

	assertDefaults := func() {
		cCnr, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Zero(t, cCnr.GetVerb())
	}

	assertBinary := func(baseAssert func()) {
		val2 := filled

		require.NoError(t, val2.Unmarshal(val.Marshal()))
		baseAssert()
	}

	assertJSON := func(baseAssert func()) {
		val2 := filled

		jd, err := val.MarshalJSON()
		require.NoError(t, err)

		require.NoError(t, val2.UnmarshalJSON(jd))
		baseAssert()
	}

	val.WriteToV2(&m)

	assertDefaults()
	assertBinary(assertDefaults)
	assertJSON(assertDefaults)

	// set value

	assertVerb := func(verb v2session.ObjectSessionVerb) {
		cCnr, ok := m.GetBody().GetContext().(*v2session.ObjectSessionContext)
		require.True(t, ok)
		require.Equal(t, verb, cCnr.GetVerb())
	}

	for from, to := range map[session.ObjectVerb]v2session.ObjectSessionVerb{
		session.VerbObjectPut:       v2session.ObjectVerbPut,
		session.VerbObjectGet:       v2session.ObjectVerbGet,
		session.VerbObjectHead:      v2session.ObjectVerbHead,
		session.VerbObjectSearch:    v2session.ObjectVerbSearch,
		session.VerbObjectRangeHash: v2session.ObjectVerbRangeHash,
		session.VerbObjectRange:     v2session.ObjectVerbRange,
		session.VerbObjectDelete:    v2session.ObjectVerbDelete,
	} {
		val.ForVerb(from)

		val.WriteToV2(&m)

		assertVerb(to)
		assertBinary(func() { assertVerb(to) })
		assertJSON(func() { assertVerb(to) })
	}
}

func TestObject_AssertVerb(t *testing.T) {
	var x session.Object

	const v1, v2 = session.VerbObjectGet, session.VerbObjectPut

	require.False(t, x.AssertVerb(v1, v2))

	x.ForVerb(v1)
	require.True(t, x.AssertVerb(v1))
	require.False(t, x.AssertVerb(v2))
	require.True(t, x.AssertVerb(v1, v2))
	require.True(t, x.AssertVerb(v2, v1))
}

func TestObject_Issuer(t *testing.T) {
	var token session.Object
	signer := randSigner()

	require.Zero(t, token.Issuer())

	require.NoError(t, token.Sign(signer))

	var issuer user.ID

	user.IDFromKey(&issuer, signer.PublicKey)

	require.True(t, token.Issuer().Equals(issuer))
}

func TestObject_Sign(t *testing.T) {
	val := sessiontest.Object()

	require.NoError(t, val.Sign(randSigner()))

	require.True(t, val.VerifySignature())
}
