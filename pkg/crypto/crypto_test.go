package crypto

import "testing"

func TestComparePasswords(t *testing.T) {
	type args struct {
		hashedPwd string
		plainPwd  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "it should return true if hash corresponds to password",
			args: args{
				hashedPwd: testStringHashSaltFromSha256,
				plainPwd:  testStringSha256,
			},
			want: true,
		},
		{
			name: "it should return false if hash does not correspond to given password",
			args: args{
				hashedPwd: testStringHashSaltFromSha256,
				plainPwd:  "a bunch of crap",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePasswords(tt.args.hashedPwd, tt.args.plainPwd); got != tt.want {
				t.Errorf("ComparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashAndSalt(t *testing.T) {
	type args struct {
		plainPwd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "it should return a valid salted hash fro the given password ",
			args:    args{plainPwd: testStringSha256},
			want:    "dynamicResult",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashAndSalt(tt.args.plainPwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashAndSalt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == "dynamicResult" {
				if !ComparePasswords(got, tt.args.plainPwd) {
					t.Errorf("HashAndSalt() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSha256Hash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "it should return the correct sha256 hash",
			args: args{s: "testString"},
			want: testStringSha256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha256Hash(tt.args.s); got != tt.want {
				t.Errorf("Sha256Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePasswordHash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "it should return false if the string is empty",
			args: args{s: ""},
			want: false,
		},
		{
			name: "it should return false if the string is 256hash from empty string",
			args: args{s: emptyStringSha256},
			want: false,
		},
		{
			name: "it should return true if the string is a valid hash",
			args: args{s: testStringSha256},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePasswordHash(tt.args.s); got != tt.want {
				t.Errorf("ValidatePasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
