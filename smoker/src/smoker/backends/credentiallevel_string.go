// Code generated by "stringer -type=CredentialLevel"; DO NOT EDIT.

package backends

import "fmt"

const _CredentialLevel_name = "UnknownUserAdmin"

var _CredentialLevel_index = [...]uint8{0, 7, 11, 16}

func (i CredentialLevel) String() string {
	if i < 0 || i >= CredentialLevel(len(_CredentialLevel_index)-1) {
		return fmt.Sprintf("CredentialLevel(%d)", i)
	}
	return _CredentialLevel_name[_CredentialLevel_index[i]:_CredentialLevel_index[i+1]]
}
