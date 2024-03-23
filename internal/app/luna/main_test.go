package luna

import "testing"

func TestIsLuna(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    bool
		wantErr bool
	}{
		{
			name:    "Good test",
			value:   "1115",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Bad test",
			value:   "1116",
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsLuna(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsLuna() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsLuna() got = %v, want %v", got, tt.want)
			}
		})
	}
}
