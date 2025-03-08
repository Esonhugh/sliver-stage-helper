package shellcoder

import (
	_ "embed"
	"fmt"

	"github.com/keystone-engine/keystone/bindings/go/keystone"
	log "github.com/sirupsen/logrus"
)

//go:embed Stage2.asm
var TemplateStr string

func GenerateLinuxX64ShellcodeFromBytes(payload []byte) ([]byte, error) {
	// 初始化 Keystone，选择架构和模式（例如 x86 32位）
	ks, err := keystone.New(keystone.ARCH_X86, keystone.MODE_64)
	if err != nil {
		log.Errorf("init failed:", err)
		return nil, err
	}
	defer ks.Close() // 确保释放资源

	// 设置汇编语法（例如 Intel 语法）
	err = ks.Option(keystone.OPT_SYNTAX, keystone.OPT_SYNTAX_INTEL)
	if err != nil {
		log.Errorf("set option failed:", err)
		return nil, err
	}

	var payload_len = len(payload)
	assembly := fmt.Sprintf(TemplateStr, fmt.Sprintf("0x%x", payload_len))
	assembly = deComment(assembly)
	// 汇编指令
	insn, _, ok := ks.Assemble(assembly, 0)
	if !ok {
		log.Fatal("asm failed, ", ks.LastError())
	}

	// 输出机器码的十六进制表示
	log.Tracef("payload_len: %v(%x)", payload_len, payload_len)
	log.Tracef("payload_prefix: %v(%x)", payload[:10], payload[:10])
	log.Tracef("payload stage2 payload: %x", insn)
	log.Tracef("len of pre-payload payload: %v", len(insn))

	return MergeBytes(insn, payload), nil
}
