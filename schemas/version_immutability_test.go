package schemas

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// publishedVersionSHA256 freezes every published JSON contract. Adding a new
// document requires adding its digest; changing or removing an existing
// document fails this test.
var publishedVersionSHA256 = map[string]string{
	"browser-registration.1.2.json":        "b3d31dc0ca833169ca34f3cc1c1517b3c86efc5c47ba2a0c9bdcb87731e65c2d",
	"browser-registration-call.1.2.json":   "175ad75d7d51d8adf8d601111dcbf220417470d08908dfdeadb7c5f1d08daf2b",
	"browser-registration.1.1.json":        "ee14fbb9ddde9bdd63ce31016cbc6f210a435694936e2d01403ba514fcd137ee",
	"browser-registration-call.1.1.json":   "a8ca1819faed5689caccac42fd35b2113c9a305e9fc93ea916e36827d4315a75",
	"browser-registration-input.1.0.json":  "9512ea3430997add675ac9c675143ad177d2d2d74fbe385f682cf29ce0e4e2c9",
	"1.0.0.json":                           "31e4b67462763807571bdd496e7d11e1e614eea43299297c7ca5301c7fa01076",
	"1.1.0.json":                           "f1652bbf473adea57bc3cc01ec2a5caf5389d39c2a4950d51857063adbd40e64",
	"1.1.1.json":                           "ffbabb57b7c6334140c18b209017a826cef5b9e607d18d263e1603e43ee323b2",
	"1.2.0.json":                           "cfd49c4d357097eb0413fac3b07d28617cfe95d7e728b75db9446e6c2b575e15",
	"1.3.0.json":                           "9099b67b0ed149f95b5413c72643dfd27734d0b7f30907023d1a4a8834a23f0e",
	"1.4.0.json":                           "7c41cf2e65c9fbecc00ca2a19d76a92d9ada8a0b2f724038aecb0497de717b1e",
	"1.5.0.json":                           "9bc90ef03383250888ec55040180a8d7e2aeb8f5d59084ad238a502f2c3a7b61",
	"1.6.0.json":                           "a447ce60caf4b12d7e142c36e125b09b9fbcc4c8610ff454207d7910df10776e",
	"1.7.0.json":                           "ca032c02f1ad51e9386a1d55d32f11d66137bd0f9c2af459e8c2a1910e4ddff8",
	"1.8.0.json":                           "ea842defb84380c6f9f6d79c3c54575c703902dd89d75322968cb80bc073e7d6",
	"1.9.0.json":                           "87b72744ea71785613cd1775968503e81312b38d2cee3a00ec5b8bf04bc976e1",
	"1.9.1.json":                           "a055a67d393dacf3c9beac30732adb59e7a81427eb8701e56e5662359c5e624d",
	"1.9.2.json":                           "ef2d9b19276beb490e400c8941649c9574d37c0c191384b8556e7f90aece18f3",
	"ansible.1.0.json":                     "c67738d98732a177863421f3edd062f98aadba325a672e9be344478dfb41c6d6",
	"browser-authentication-call.1.0.json": "586af5315001334ba7ceb69048f31b278ef16b0984ab05da6b7b361a2e035672",
	"browser-authentication-call.1.1.json": "361aa798cefb4a172ecfda8c794a970faeb05919557a73bacb143af57332d35a",
	"browser-registration-call.1.0.json":   "b311aeefb1b6c2b8d675a0839c654140105223f53994eb14260b9387aa293a38",
	"browser-registration.1.0.json":        "6613ba3eecf0d073554cc6e6a4a56ebf2e9e3264714cd49fd740e6e784044605",
	"browser-authentication.1.0.json":      "8ccf16281a83783d0b342312edb0b737a53129fca0944d8f00f31acc91f461ad",
	"browser-authentication.1.1.json":      "61ba61569e062959213533cf95363cfe07b226241e75e5c4198d690d35debbd2",
	"browser.1.5.json":                     "6dda6191ca4d9899183f29c170e80e7c295ecb33f3e32c8994898e62eb8fa4c8",
	"browser.1.6.json":                     "396d36fff165b2bf4fd6ada45cacad7365f330ab3ae16fc95d9a244e34f819bc",
	"browser.1.7.json":                     "feed4f71655b232fe6a87db285ac686615f9e2076c64b751e5c945c9463f446b",
	"runtime.1.0.json":                     "c8ed61ae855c828767a30d94e667bd7f0b3bed75ee8e36407f815d789fe6cd31",
}

// publishedMarkdownSHA256 freezes every published Markdown document except
// CHANGELOG.md, which remains the mutable release ledger.
var publishedMarkdownSHA256 = map[string]string{
	"1.0.0.md":                           "b4697ee580838bc7e6e0722d01219534086b8e9368b406c1ccae1c2440dad671",
	"1.1.0.md":                           "f38d92ce04684f3aad9ee9012fb6ef1b7a1d1972e6af668ec98bb459070993f7",
	"1.1.1.md":                           "57577c86fbaab8fdfa362be6dd6b8c01ceaa271084ed6e836e123081032e095e",
	"1.2.0.md":                           "e12596201e6ba35cc32b5cb629e2d5ff0702ebd046bd7d96f01bb47d685fb34c",
	"1.3.0.md":                           "13a71b79e3fd2868968e6fce6d31c73995ac7c549dcf20f94d72af01fa2b9486",
	"1.4.0.md":                           "21c0c3d8461b62b0cf18326c83306ff354d9c4e79720e4a483a1eb0e224f170d",
	"1.5.0.md":                           "cae9bf3d5a830031d7c85297e031450c809c44983baaf6b34d49b28a5c8761bc",
	"1.6.0.md":                           "72a458723674c95765b416ed0a3efd82e46cb9dc585f5530887e9fa682d07615",
	"1.7.0.md":                           "0d040e7cd390893d5999425fbe3dcb190ff39ff0b0b44eb41e42e79582c9bdfb",
	"1.8.0.md":                           "cea8f392040ef4cc99f77869cf4875fba484a28027677dfee8b9b0824b59531d",
	"1.9.0.md":                           "9e594f3326811947799805580d6a0f4e508b4ef880187171311fe0b34c69b98d",
	"1.9.1.md":                           "de2ba4fe07aa22ecfcdfbfe243caf09e950a5592b3d98d9cf826a0bb5cf998f3",
	"1.9.2.md":                           "92171354a3fe60faa468fe0e43f36733f535ab6a3de8b57b82afe245bf6a57dc",
	"ansible.1.0.md":                     "fc7ac843633004e33bac0bee34e178a4de894491a222162875c6f6f0bfcea504",
	"arazzo.md":                          "f4591078579df231f030d8877667fd2c2b50c8e713508361d516f67a6ed67490",
	"article.md":                         "2a1b30f5515efa95e906168033f3f8c305df9b34a026e32dabf60397ebb4bcda",
	"browser-authentication-call.1.0.md": "c1b3ace683f2bb301c25dff6ed0b49bc14f2e8ba6e964c7855c76f4a736ceb38",
	"browser-authentication-call.1.1.md": "b0dbef259c3a8c2e8ea82b2cb741d62d781cf0aae60004b981bd465f12a81800",
	"browser-authentication.1.0.md":      "4747e20a4f73bbb1c6205e822bf5dd188faee7eab31144ac45cb654679d70a74",
	"browser-authentication.1.1.md":      "182691589ca39d28c983a79b6ad73f7b21fb7b0e9f15c0ac09352e70b451f128",
	"browser-registration-call.1.0.md":   "18743f0817580466721c16208eb85b9e617e495ef3458cfbc585f25ecac97f84",
	"browser-registration-call.1.1.md":   "a2d7b101ac5134b7f7165737f3f3b9db4afa8f18784b39d4b11cc7ef3a9eb3af",
	"browser-registration-call.1.2.md":   "668b85d37b14e171d2197c3612d2373ec4deae8940f9d7b8194788c9fc9fa419",
	"browser-registration-input.1.0.md":  "35d59168f4041e1a20406a655c3f7a381018afaac78467bae31e657fefaae735",
	"browser-registration.1.0.md":        "089eae4bc81ec6e6dfa8d1b263218624f9d31bdbc236032570bf63a29e4db53c",
	"browser-registration.1.1.md":        "ee1fb0842191f1df08276a2c751478509903e277db91a89cbf14de7461a4a97c",
	"browser-registration.1.2.md":        "94e16ced40984a32a308bef7f9c13d40ea1b7c627c3519d49c3930d509887154",
	"browser.1.5.md":                     "05be7331270084b4daa373632ef75b2db357c49836a376b8ab9b84444fb0d632",
	"browser.1.6.md":                     "1e79e778060e458c078f3f3e65c8a87484887b144222dd3221e8b3262da5c0f1",
	"browser.1.7.md":                     "76ff87cab25bee91300f38414c542d31c3d9a23a63eda8443fcec702b4785d83",
	"runtime.1.0.md":                     "9f77d78f250d8a1e98a1e9c0ebc0536cd8091250f03701a53be2b0d7280a5576",
}

func TestPublishedVersionDocumentsAreImmutable(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		assertPublishedVersionDocumentsImmutable(t, "*.json", "JSON", "", publishedVersionSHA256)
	})
	t.Run("Markdown", func(t *testing.T) {
		assertPublishedVersionDocumentsImmutable(t, "*.md", "Markdown", "CHANGELOG.md", publishedMarkdownSHA256)
	})
}

func assertPublishedVersionDocumentsImmutable(t *testing.T, pattern, kind, excludedName string, manifest map[string]string) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join("..", "versions", pattern))
	if err != nil {
		t.Fatal(err)
	}
	actualNames := make([]string, 0, len(paths))
	for _, path := range paths {
		name := filepath.Base(path)
		if name == excludedName {
			continue
		}
		actualNames = append(actualNames, name)
		want, ok := manifest[name]
		if !ok {
			t.Errorf("published %s document %s is missing from the SHA-256 manifest", kind, name)
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Error(err)
			continue
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != want {
			t.Errorf("published %s document %s changed: SHA-256 = %s, want %s", kind, name, got, want)
		}
	}
	sort.Strings(actualNames)
	wantNames := make([]string, 0, len(manifest))
	for name := range manifest {
		wantNames = append(wantNames, name)
	}
	sort.Strings(wantNames)
	if len(actualNames) != len(wantNames) {
		t.Fatalf("versions/%s membership = %v, manifest membership = %v", pattern, actualNames, wantNames)
	}
	for i := range actualNames {
		if actualNames[i] != wantNames[i] {
			t.Fatalf("versions/%s membership = %v, manifest membership = %v", pattern, actualNames, wantNames)
		}
	}
}
