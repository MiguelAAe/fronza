package addresssvc

import "testing"

func TestCheckPostCode(t *testing.T) {
	tcs := []struct {
		name     string
		postcode string
		want     bool
	}{
		{
			name:     "happy path",
			postcode: "e149ed",
			want:     true,
		},
		{
			name:     "happy path",
			postcode: " e149ed ",
			want:     true,
		},
		{
			name:     "happy path",
			postcode: "e14 9 ed",
			want:     true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := CheckPostCode(tc.postcode)
			if got != tc.want {
				t.Errorf("did not match expected value")
			}
		})
	}
}

func TestGetPostcodeCoordinates(t *testing.T) {
	tcs := []struct {
		name          string
		postcode      string
		wantLatitude  string
		wantLongigute string
	}{
		{
			name:          "happy path: present postcode",
			postcode:      "n168rp",
			wantLatitude:  "51.556324",
			wantLongigute: "-0.080866",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			gotLatitude, gotLongitude := GetPostcodeCoordinates(tc.postcode)
			if gotLatitude != tc.wantLatitude {
				t.Errorf("latitude did not match expected value, got: %s, want: %s", gotLatitude, tc.wantLatitude)
			}

			if gotLongitude != tc.wantLongigute {
				t.Errorf("longitude did not match expected value, got: %s, want: %s", gotLongitude, tc.wantLongigute)
			}
		})
	}
}
