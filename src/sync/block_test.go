/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"encoding/json"
	"os"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var BlockJson = "{\"block\":{\"header\":{\"block_version\":2030200,\"chain_id\":\"chain1\",\"block_height\":30,\"block_hash\":\"Ag42/wT1z9OBZ956lzM+GMMsz6eWnY7Xe1Td5/SPCl4=\",\"pre_block_hash\":\"pdwLSpdM8FXvA7PDJXC4otrRAdwi3yfhHMB9gd0aMB0=\",\"pre_conf_height\":22,\"tx_count\":1,\"tx_root\":\"QKXkc0pPzhlTDpZz48NHC6GidGEaV0rrZP1V3gftt4k=\",\"dag_hash\":\"CNp8RcsgQ3fn5CJJzaVxP6hlEW3btMtaGUmy5bQ4pqs=\",\"rw_set_root\":\"f9dibAngnnhjOF9mWwoyW/HWhFU3Ie+/Ug0wymh3k8o=\",\"block_timestamp\":1702460320,\"proposer\":{\"org_id\":\"wx-org3.chainmaker.org\",\"member_info\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNmakNDQWlTZ0F3SUJBZ0lEQXhYWU1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTXk1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2daY3hDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jekxtTm9ZV2x1YldGclpYSXViM0puCk1SSXdFQVlEVlFRTEV3bGpiMjV6Wlc1emRYTXhMekF0QmdOVkJBTVRKbU52Ym5ObGJuTjFjekV1YzJsbmJpNTMKZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRQpEdHpTMnhZVFVmT1l1dlFxK2MzS2tMMzZaVFkrTkNIYXJ0VUhTOGNDVzU0S1FYUmZEUVJDMGU5L0pSMUk3a3A5CjlneTlmc2h2bUVTdy9nUXlrZktYOXFOcU1HZ3dEZ1lEVlIwUEFRSC9CQVFEQWdiQU1Da0dBMVVkRGdRaUJDRGwKZXdHZXBTRS9kcGRjcStuR0NIZUMwejdRcEI2Zlh5WnM2NFR1ekJrL1VqQXJCZ05WSFNNRUpEQWlnQ0RHbzVRYwpMd1lJdVVGMDN3RWEwM29wMHRPdGVBNFloc3ZVdVpFb3ZwSlhHVEFLQmdncWhrak9QUVFEQWdOSUFEQkZBaUJsCnU0ZExObGxoOTFSNWprTE9JMklXY1BkNGh0MWpUaC96Z3Q4TVVFWkY2QUloQUp6TTU2azdiY1JmcUF3RGVDREIKekdGMlQxTktVSEZucUp1NlliRFg4RDVsCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\"},\"signature\":\"MEQCIFBzMeUG/A67UaVqQUParBgMBimG+bKd9AtesPRnqlWuAiBkTrgev/YdVlmUPg20r4RGdDpzuafib+IOvUTW93zr/w==\"},\"dag\":{\"vertexes\":[{}]},\"txs\":[{\"payload\":{\"chain_id\":\"chain1\",\"tx_id\":\"17a05aa42d8d673cca4160edc45863596289046afae0442c9a5bcd6bc5347359\",\"timestamp\":1702460320,\"contract_name\":\"CONTRACT_MANAGE\",\"method\":\"INIT_CONTRACT\",\"parameters\":[{\"key\":\"data\"},{\"key\":\"CONTRACT_NAME\",\"value\":\"RVZNX0VSQzIwXzE=\"},{\"key\":\"CONTRACT_VERSION\",\"value\":\"MS4w\"},{\"key\":\"CONTRACT_RUNTIME_TYPE\",\"value\":\"RVZN\"},{\"key\":\"CONTRACT_BYTECODE\",\"value\":\"YGBgQFJrAzsuPJ/QgDzoAAAAYANVNBVhAB9XYACA/VtgA1RgAIAzc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIIGQVVAzc///////////////////////////FmAAc///////////////////////////Fn/d8lKtG+LIm2nCsGj8N42qlSun8WPEoRYo9VpN9SOz72ADVGBAUYCCgVJgIAGRUFBgQFGAkQOQo2ELa4BhANtgADlgAPMAYGBgQFJgBDYQYQCZV2AANXwBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJAEY/////8WgGMG/d4DFGEAnleAYwlep7MUYQEsV4BjGBYN3RRhAYZXgGMjuHLdFGEBr1eAYzE85WcUYQIoV4BjcKCCMRRhAldXgGOV2JtBFGECpFeAY6kFnLsUYQMyV4Bj3WLtPhRhA4xXW2AAgP1bNBVhAKlXYACA/VthALFhA/hWW2BAUYCAYCABgoEDglKDgYFRgVJgIAGRUIBRkGAgAZCAg4NgAFuDgRAVYQDxV4CCAVGBhAFSYCCBAZBQYQDWVltQUFBQkFCQgQGQYB8WgBVhAR5XgIIDgFFgAYNgIANhAQAKAxkWgVJgIAGRUFtQklBQUGBAUYCRA5DzWzQVYQE3V2AAgP1bYQFsYASAgDVz//////////////////////////8WkGAgAZCRkIA1kGAgAZCRkFBQYQQxVltgQFGAghUVFRWBUmAgAZFQUGBAUYCRA5DzWzQVYQGRV2AAgP1bYQGZYQUjVltgQFGAgoFSYCABkVBQYEBRgJEDkPNbNBVhAbpXYACA/VthAg5gBICANXP//////////////////////////xaQYCABkJGQgDVz//////////////////////////8WkGAgAZCRkIA1kGAgAZCRkFBQYQUpVltgQFGAghUVFRWBUmAgAZFQUGBAUYCRA5DzWzQVYQIzV2AAgP1bYQI7YQhOVltgQFGAgmD/FmD/FoFSYCABkVBQYEBRgJEDkPNbNBVhAmJXYACA/VthAo5gBICANXP//////////////////////////xaQYCABkJGQUFBhCFNWW2BAUYCCgVJgIAGRUFBgQFGAkQOQ81s0FWECr1dgAID9W2ECt2EIm1ZbYEBRgIBgIAGCgQOCUoOBgVGBUmAgAZFQgFGQYCABkICDg2AAW4OBEBVhAvdXgIIBUYGEAVJgIIEBkFBhAtxWW1BQUFCQUJCBAZBgHxaAFWEDJFeAggOAUWABg2AgA2EBAAoDGRaBUmAgAZFQW1CSUFBQYEBRgJEDkPNbNBVhAz1XYACA/VthA3JgBICANXP//////////////////////////xaQYCABkJGQgDWQYCABkJGQUFBhCNRWW2BAUYCCFRUVFYFSYCABkVBQYEBRgJEDkPNbNBVhA5dXYACA/VthA+JgBICANXP//////////////////////////xaQYCABkJGQgDVz//////////////////////////8WkGAgAZCRkFBQYQq4VltgQFGAgoFSYCABkVBQYEBRgJEDkPNbYECAUZCBAWBAUoBgE4FSYCABf1Rlc3RDaGFpbk1ha2VyVG9rZW4AAAAAAAAAAAAAAAAAgVJQgVZbYACBYAFgADNz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYACFc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIIGQVVCCc///////////////////////////FjNz//////////////////////////8Wf4xb4eXr7H1b0U9xQn0ehPPdAxTA97IpHlsgCsjHw7klhGBAUYCCgVJgIAGRUFBgQFGAkQOQo2ABkFCSkVBQVltgA1SBVltgAIBgAWAAhnP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBgADNz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVJBQgmAAgIdz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVBAVgBVhBflXUIKBEBVbgBVhBoNXUGAAgIVz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVINgAICHc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIFQBEBVbFWEIQVeCYACAhnP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBgAIKCVAGSUFCBkFVQgmAAgIdz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYACCglQDklBQgZBVUH///////////////////////////////////////////4EQFWEH01eCYAFgAIdz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYAAzc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIGAAgoJUA5JQUIGQVVBbg3P//////////////////////////xaFc///////////////////////////Fn/d8lKtG+LIm2nCsGj8N42qlSun8WPEoRYo9VpN9SOz74VgQFGAgoFSYCABkVBQYEBRgJEDkKNgAZFQYQhGVltgAJFQW1CTklBQUFZbYBKBVltgAIBgAINz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVJBQkZBQVltgQIBRkIEBYEBSgGADgVJgIAF/VENNAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACBUlCBVltgAIFgAIAzc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIFQQFYAVYQmiV1BgAICEc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIFSCYACAhnP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBUARAVWxVhCq1XgWAAgDNz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYACCglQDklBQgZBVUIFgAICFc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIGAAgoJUAZJQUIGQVVCCc///////////////////////////FjNz//////////////////////////8Wf93yUq0b4sibacKwaPw3jaqVK6fxY8ShFij1Wk31I7PvhGBAUYCCgVJgIAGRUFBgQFGAkQOQo2ABkFBhCrJWW2AAkFBbkpFQUFZbYABgAWAAhHP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBgAINz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVJBQkpFQUFYAoWVienpyMFggJZE/R6lfV/obNOchph0QuRIBwjsDd47+Oai3P2pV+ccAKQ==\"}],\"limit\":{\"gas_limit\":130000}},\"sender\":{\"signer\":{\"org_id\":\"wx-org1.chainmaker.org\",\"member_type\":1,\"member_info\":\"LK4/KplYsQcFU2All0UxorspVALdt/tgHuZ8QxiME2M=\"},\"signature\":\"MEUCIQCrui2zf66Ibk40eVTU8HeZ1f0QpMqrfXqpyz7ivzbiyQIgZ4N3fswmu2dlXD6MGW4aUxwbVdeBFGNe2O1RDE00+hM=\"},\"endorsers\":[{\"signer\":{\"org_id\":\"wx-org1.chainmaker.org\",\"member_info\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNkVENDQWh5Z0F3SUJBZ0lEQ1cvek1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTVM1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3hMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2dZOHhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jeExtTm9ZV2x1YldGclpYSXViM0puCk1RNHdEQVlEVlFRTEV3VmhaRzFwYmpFck1Da0dBMVVFQXhNaVlXUnRhVzR4TG5OcFoyNHVkM2d0YjNKbk1TNWoKYUdGcGJtMWhhMlZ5TG05eVp6QlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJJY09uaThyWUlyVwpLV3E1aS93QzcrNU94cHdoeDBmUGQrMjcwSG9vWHNMTkEzL2tLN0N3QThLMzhIeFBGZkFIaU85aE5zaHdGK2k3CnhZTzdkaVJVcGJxamFqQm9NQTRHQTFVZER3RUIvd1FFQXdJR3dEQXBCZ05WSFE0RUlnUWdmczUxQlNnSHZrNE0KaURjQ0pQZ2h4WTNUR1BVR0lXdVZXb1pEeUMzN1l0c3dLd1lEVlIwakJDUXdJb0FnRDBhdnY5Y0VENC9BSkhsUwpNcGZBK290SFlLcjhjdVhyZ3ZneGVWbGJsRGN3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnRTEwTnMxZndDbjFVCmZuREpuZjBTVEY1SXBtMzZ5c2NwdmJCbDBMUG91UUVDSUJzNFVQR3NVL2dYbHZ6Q0dqWW5DRkdUZmo0bnkrRjIKdXNwdkxya0ZvemVpCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K\"},\"signature\":\"MEUCIEWS9pitu7ioiCtWyGdfMmT9d0VwRDuONtoBaqB69ottAiEAzfbCHDixvXYD/KKeA8J/5VWjGbij4ZSWMNv5MtiOR1E=\"},{\"signer\":{\"org_id\":\"wx-org2.chainmaker.org\",\"member_info\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNkakNDQWh5Z0F3SUJBZ0lEQzdGMk1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTWk1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3lMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2dZOHhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jeUxtTm9ZV2x1YldGclpYSXViM0puCk1RNHdEQVlEVlFRTEV3VmhaRzFwYmpFck1Da0dBMVVFQXhNaVlXUnRhVzR4TG5OcFoyNHVkM2d0YjNKbk1pNWoKYUdGcGJtMWhhMlZ5TG05eVp6QlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJFQ3FtKzEyTG1SUApmZHdUMjRzWittcEp4Y1l3bzI2cmsyWTRSUGpsdFV1KzNGNEFEell1ZXBzc1VBR2t2MHlvT00vRUhyMUNRUUdWClh3TS9OVjdvMnIramFqQm9NQTRHQTFVZER3RUIvd1FFQXdJR3dEQXBCZ05WSFE0RUlnUWdzMGZzQms1Q3hBN2MKVEppbUZadVZ5YUJiNTlRcUQwZi9oOGpQUmZPaDAvd3dLd1lEVlIwakJDUXdJb0FnZGVteWhaSUNaSXFJYkR5ZQpuVW5FK0hKOXBrTDB3QjY5MzJWcXBlY081TlF3Q2dZSUtvWkl6ajBFQXdJRFNBQXdSUUloQUxWNnBjMFZJTUNmCnI5dFNTZ2E4aWp6b0VCdEVGSzNZWUlaWXpSb1dTY0p1QWlBeHErR3N0aUtveXZqSkdadkFYMWZWWmNlTXRmbEEKTEpnUmNsZE1CSlZnYXc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\"},\"signature\":\"MEUCIQCMbDGNH/WLa2B7i3qWa/+SWZO/uIMUdxE56KCBfRW5gwIgRWjXzZ6TkCrhR2Psvpa3upoVgE39fQn6r5mBDdZv6u8=\"},{\"signer\":{\"org_id\":\"wx-org3.chainmaker.org\",\"member_info\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNkakNDQWh5Z0F3SUJBZ0lEQlNvNU1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTXk1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2dZOHhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jekxtTm9ZV2x1YldGclpYSXViM0puCk1RNHdEQVlEVlFRTEV3VmhaRzFwYmpFck1Da0dBMVVFQXhNaVlXUnRhVzR4TG5OcFoyNHVkM2d0YjNKbk15NWoKYUdGcGJtMWhhMlZ5TG05eVp6QlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJBM2pEZCs1d0N5Sgp2WW05Vko3eWJMU0R6STVzS3p2YXZOYkpvNVZJNU5GaERwZlpyTm1qN29hWlgzOTdjUkgzWWFvZnhlSnZ6ZnhKCnRHRDQxeEk3Vk9XamFqQm9NQTRHQTFVZER3RUIvd1FFQXdJR3dEQXBCZ05WSFE0RUlnUWcycklJNExwTDRmS1kKVGVORURrcWpMZWJ4M0k0U1lRMWtBbTJYM0JmNmpRSXdLd1lEVlIwakJDUXdJb0FneHFPVUhDOEdDTGxCZE44QgpHdE42S2RMVHJYZ09HSWJMMUxtUktMNlNWeGt3Q2dZSUtvWkl6ajBFQXdJRFNBQXdSUUloQVB4ZUR3Yis5Z2ZwCmdMdWgyWUJrV0NqbGVNdmg4NFNVUlcvN290QXFlOTVsQWlBeDMxdVUzbFZiUXJSWlhISmZldng3cVZ0am1ja3YKWmFDc1pnUWJhbmZueHc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\"},\"signature\":\"MEYCIQCFRwSgpgLPbEBXKG7QNaS8f9xOHr7Fcf6utdoeIgq7RQIhAOU7KXZ6+Q3BbbuJTPKt2wgJfdK7agZCbfKZPD67pon6\"}],\"result\":{\"contract_result\":{\"result\":\"CgtFVk1fRVJDMjBfMRIDMS4wGAUqqwEKFnd4LW9yZzEuY2hhaW5tYWtlci5vcmcQARogLK4/KplYsQcFU2All0UxorspVALdt/tgHuZ8QxiME2MiI2NsaWVudDEuc2lnbi53eC1vcmcxLmNoYWlubWFrZXIub3JnKgZDTElFTlQyQDI4ZTI1YTYzMGRiZTY5MjYwMDZjNmMwYjJlZjYyYWI0MDUyZTU2Y2YwZWMxZmY3ZmRkY2QyZTIzY2Q5YjQ2NGEyKDg4MmU1NDBmNDhkNmFlYjljNjRmOWFmZmM3NDg2MjRhOTU3ZDNkNGY=\",\"gas_used\":29853,\"contract_event\":[{\"topic\":\"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"tx_id\":\"17a05aa42d8d673cca4160edc45863596289046afae0442c9a5bcd6bc5347359\",\"contract_name\":\"EVM_ERC20_1\",\"contract_version\":\"1.0\",\"event_data\":[\"0000000000000000000000000000000000000000000000000000000000000000\",\"000000000000000000000000171262347a59fded92021a32421a5dad05424e03\",\"0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000\"]}]},\"rw_set_hash\":\"f9dibAngnnhjOF9mWwoyW/HWhFU3Ie+/Ug0wymh3k8o=\"}}],\"additional_data\":{\"extra_data\":{\"TBFTAddtionalDataKey\":\"CAEQHhgMIAMqIAIONv8E9c/TgWfeepczPhjDLM+nlp2O13tU3ef0jwpeMpUJCi5RbVVqQzI5VkVCdFJjRDJqUXY4a1pMbUpucU5lRWp3alNKS2UxN3pLdXlzVXNKEuIICAESLlFtVWpDMjlWRUJ0UmNEMmpRdjhrWkxtSm5xTmVFandqU0pLZTE3ekt1eXNVc0oYHiAMKiACDjb/BPXP04Fn3nqXMz4YwyzPp5adjtd7VN3n9I8KXjKHCAq7BwoWd3gtb3JnNC5jaGFpbm1ha2VyLm9yZxqgBy0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDZnpDQ0FpU2dBd0lCQWdJREFSRTlNQW9HQ0NxR1NNNDlCQU1DTUlHS01Rc3dDUVlEVlFRR0V3SkRUakVRCk1BNEdBMVVFQ0JNSFFtVnBhbWx1WnpFUU1BNEdBMVVFQnhNSFFtVnBhbWx1WnpFZk1CMEdBMVVFQ2hNV2QzZ3QKYjNKbk5DNWphR0ZwYm0xaGEyVnlMbTl5WnpFU01CQUdBMVVFQ3hNSmNtOXZkQzFqWlhKME1TSXdJQVlEVlFRRApFeGxqWVM1M2VDMXZjbWMwTG1Ob1lXbHViV0ZyWlhJdWIzSm5NQjRYRFRJek1USXdNVEE0TkRNeE5Gb1hEVEk0Ck1URXlPVEE0TkRNeE5Gb3dnWmN4Q3pBSkJnTlZCQVlUQWtOT01SQXdEZ1lEVlFRSUV3ZENaV2xxYVc1bk1SQXcKRGdZRFZRUUhFd2RDWldscWFXNW5NUjh3SFFZRFZRUUtFeFozZUMxdmNtYzBMbU5vWVdsdWJXRnJaWEl1YjNKbgpNUkl3RUFZRFZRUUxFd2xqYjI1elpXNXpkWE14THpBdEJnTlZCQU1USm1OdmJuTmxibk4xY3pFdWMybG5iaTUzCmVDMXZjbWMwTG1Ob1lXbHViV0ZyWlhJdWIzSm5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUUKYjhoWmRTaERISytrcStKNkRKcVRWdm5yMTdONU4vam5tMTR0cjh5YVdmWVNoNlZ3SzEzckV1blpFeGtxdmNHegpVcSs2aFZOaEI1dEtkczdIZk5ESm5LTnFNR2d3RGdZRFZSMFBBUUgvQkFRREFnYkFNQ2tHQTFVZERnUWlCQ0JKCkkxQ0REcCtyVy9RUFVFeXJMMVh6YXlIZFgvS2ZDSUpreGxFNjRtOXFSREFyQmdOVkhTTUVKREFpZ0NEckFoaXcKdVkyNldERzRzRGpTRHowNXdtZnVYSUI0V0hmT3dWUm9UUHhnNERBS0JnZ3Foa2pPUFFRREFnTkpBREJHQWlFQQpnUksxL21UQ1dxM2hkd2lKakkvT0gxVWNPdWwraUNRN3Y0S3BvSlhmdGFnQ0lRRC9VQmRDMjRCb3FHNmE0a2FPCmdjcmlsaFVXTSs5OWREbUhGOUdpSndjNjFnPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQoSRzBFAiAzDI7qCHs2zWbhwPrB5QYEqZOXxQ8LLUlv7fneYcO/1wIhAL6DmeTjm82uHOyacH3kFwUybSfnFs7YtxSU6/c1j4QPMpEJCi5RbVRLd3lzcFQ1NUpjcThpZW4xSk15eGNrTjZBbXFnQnBDejJ2VVp5cDROUmRYEt4ICAESLlFtVEt3eXNwVDU1SmNxOGllbjFKTXl4Y2tONkFtcWdCcEN6MnZVWnlwNE5SZFgYHiAMKiACDjb/BPXP04Fn3nqXMz4YwyzPp5adjtd7VN3n9I8KXjKDCAq3BwoWd3gtb3JnMi5jaGFpbm1ha2VyLm9yZxqcBy0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDZmpDQ0FpU2dBd0lCQWdJREFxQm9NQW9HQ0NxR1NNNDlCQU1DTUlHS01Rc3dDUVlEVlFRR0V3SkRUakVRCk1BNEdBMVVFQ0JNSFFtVnBhbWx1WnpFUU1BNEdBMVVFQnhNSFFtVnBhbWx1WnpFZk1CMEdBMVVFQ2hNV2QzZ3QKYjNKbk1pNWphR0ZwYm0xaGEyVnlMbTl5WnpFU01CQUdBMVVFQ3hNSmNtOXZkQzFqWlhKME1TSXdJQVlEVlFRRApFeGxqWVM1M2VDMXZjbWN5TG1Ob1lXbHViV0ZyWlhJdWIzSm5NQjRYRFRJek1USXdNVEE0TkRNeE5Gb1hEVEk0Ck1URXlPVEE0TkRNeE5Gb3dnWmN4Q3pBSkJnTlZCQVlUQWtOT01SQXdEZ1lEVlFRSUV3ZENaV2xxYVc1bk1SQXcKRGdZRFZRUUhFd2RDWldscWFXNW5NUjh3SFFZRFZRUUtFeFozZUMxdmNtY3lMbU5vWVdsdWJXRnJaWEl1YjNKbgpNUkl3RUFZRFZRUUxFd2xqYjI1elpXNXpkWE14THpBdEJnTlZCQU1USm1OdmJuTmxibk4xY3pFdWMybG5iaTUzCmVDMXZjbWN5TG1Ob1lXbHViV0ZyWlhJdWIzSm5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUUKK0R5dEJOZXZFSjRYWEQvdFpxdDJCWTcrK3NrQ0NqMHRDajA3T0kyRmZVY1VJTTliTVBGclloMEVEZ0pCUy9LSgowTm5lRWVPNG8rdkFCd3B4N2xEQ2VhTnFNR2d3RGdZRFZSMFBBUUgvQkFRREFnYkFNQ2tHQTFVZERnUWlCQ0FFCmdhVkU3T3FHSCs4YzJHc3p0bGkxdmxkazZsWEk0MnIwUzBZUGNNT3hWVEFyQmdOVkhTTUVKREFpZ0NCMTZiS0YKa2dKa2lvaHNQSjZkU2NUNGNuMm1RdlRBSHIzZlpXcWw1dzdrMURBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlFQQprdm8rMjhsQzFHcWdJck5OYVh5VmdXcWdSQkFCYnUrZVV3YTBwSGVsTkdvQ0lCU09sRi9tb3JEUTBJcnU4ZGVnCkJHRzFHRXRuYWx4WUVOdm1XMnpWUDdmdgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tChJHMEUCIDIdWc2mERyF0FVMMxproVxPL23czaqdyRv2UTiAbrRAAiEA5KxxqW+7exwaiX1gwxAZUNtfKhKAdMZm/nZ3fBUf3tcykQkKLlFtWkdLUHd1M0s5dUNkbzN4NEVzUjVONG5tUjRxcGs3UTU2Ukx6c05kSGJLdzgS3ggIARIuUW1aR0tQd3UzSzl1Q2RvM3g0RXNSNU40bm1SNHFwazdRNTZSTHpzTmRIYkt3OBgeIAwqIAIONv8E9c/TgWfeepczPhjDLM+nlp2O13tU3ef0jwpeMoMICrcHChZ3eC1vcmczLmNoYWlubWFrZXIub3JnGpwHLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNmakNDQWlTZ0F3SUJBZ0lEQXhYWU1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTXk1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2daY3hDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jekxtTm9ZV2x1YldGclpYSXViM0puCk1SSXdFQVlEVlFRTEV3bGpiMjV6Wlc1emRYTXhMekF0QmdOVkJBTVRKbU52Ym5ObGJuTjFjekV1YzJsbmJpNTMKZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRQpEdHpTMnhZVFVmT1l1dlFxK2MzS2tMMzZaVFkrTkNIYXJ0VUhTOGNDVzU0S1FYUmZEUVJDMGU5L0pSMUk3a3A5CjlneTlmc2h2bUVTdy9nUXlrZktYOXFOcU1HZ3dEZ1lEVlIwUEFRSC9CQVFEQWdiQU1Da0dBMVVkRGdRaUJDRGwKZXdHZXBTRS9kcGRjcStuR0NIZUMwejdRcEI2Zlh5WnM2NFR1ekJrL1VqQXJCZ05WSFNNRUpEQWlnQ0RHbzVRYwpMd1lJdVVGMDN3RWEwM29wMHRPdGVBNFloc3ZVdVpFb3ZwSlhHVEFLQmdncWhrak9QUVFEQWdOSUFEQkZBaUJsCnU0ZExObGxoOTFSNWprTE9JMklXY1BkNGh0MWpUaC96Z3Q4TVVFWkY2QUloQUp6TTU2azdiY1JmcUF3RGVDREIKekdGMlQxTktVSEZucUp1NlliRFg4RDVsCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkcwRQIgaBZ5QBktg17dttD7StHoBOkRwA0iub8eG6DCJbK4iVcCIQCA2I2513CoDtabPcDr0WUzfrGWVcvHWF47FThK4ucs/DrzGwosQWc0Mi93VDF6OU9CWjk1Nmx6TStHTU1zejZlV25ZN1hlMVRkNS9TUENsND0SwhsKkQkKLlFtWkdLUHd1M0s5dUNkbzN4NEVzUjVONG5tUjRxcGs3UTU2Ukx6c05kSGJLdzgS3ggIARIuUW1aR0tQd3UzSzl1Q2RvM3g0RXNSNU40bm1SNHFwazdRNTZSTHpzTmRIYkt3OBgeIAwqIAIONv8E9c/TgWfeepczPhjDLM+nlp2O13tU3ef0jwpeMoMICrcHChZ3eC1vcmczLmNoYWlubWFrZXIub3JnGpwHLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNmakNDQWlTZ0F3SUJBZ0lEQXhYWU1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTXk1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2daY3hDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jekxtTm9ZV2x1YldGclpYSXViM0puCk1SSXdFQVlEVlFRTEV3bGpiMjV6Wlc1emRYTXhMekF0QmdOVkJBTVRKbU52Ym5ObGJuTjFjekV1YzJsbmJpNTMKZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRQpEdHpTMnhZVFVmT1l1dlFxK2MzS2tMMzZaVFkrTkNIYXJ0VUhTOGNDVzU0S1FYUmZEUVJDMGU5L0pSMUk3a3A5CjlneTlmc2h2bUVTdy9nUXlrZktYOXFOcU1HZ3dEZ1lEVlIwUEFRSC9CQVFEQWdiQU1Da0dBMVVkRGdRaUJDRGwKZXdHZXBTRS9kcGRjcStuR0NIZUMwejdRcEI2Zlh5WnM2NFR1ekJrL1VqQXJCZ05WSFNNRUpEQWlnQ0RHbzVRYwpMd1lJdVVGMDN3RWEwM29wMHRPdGVBNFloc3ZVdVpFb3ZwSlhHVEFLQmdncWhrak9QUVFEQWdOSUFEQkZBaUJsCnU0ZExObGxoOTFSNWprTE9JMklXY1BkNGh0MWpUaC96Z3Q4TVVFWkY2QUloQUp6TTU2azdiY1JmcUF3RGVDREIKekdGMlQxTktVSEZucUp1NlliRFg4RDVsCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkcwRQIgaBZ5QBktg17dttD7StHoBOkRwA0iub8eG6DCJbK4iVcCIQCA2I2513CoDtabPcDr0WUzfrGWVcvHWF47FThK4ucs/AqVCQouUW1VakMyOVZFQnRSY0QyalF2OGtaTG1KbnFOZUVqd2pTSktlMTd6S3V5c1VzShLiCAgBEi5RbVVqQzI5VkVCdFJjRDJqUXY4a1pMbUpucU5lRWp3alNKS2UxN3pLdXlzVXNKGB4gDCogAg42/wT1z9OBZ956lzM+GMMsz6eWnY7Xe1Td5/SPCl4yhwgKuwcKFnd4LW9yZzQuY2hhaW5tYWtlci5vcmcaoActLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ2Z6Q0NBaVNnQXdJQkFnSURBUkU5TUFvR0NDcUdTTTQ5QkFNQ01JR0tNUXN3Q1FZRFZRUUdFd0pEVGpFUQpNQTRHQTFVRUNCTUhRbVZwYW1sdVp6RVFNQTRHQTFVRUJ4TUhRbVZwYW1sdVp6RWZNQjBHQTFVRUNoTVdkM2d0CmIzSm5OQzVqYUdGcGJtMWhhMlZ5TG05eVp6RVNNQkFHQTFVRUN4TUpjbTl2ZEMxalpYSjBNU0l3SUFZRFZRUUQKRXhsallTNTNlQzF2Y21jMExtTm9ZV2x1YldGclpYSXViM0puTUI0WERUSXpNVEl3TVRBNE5ETXhORm9YRFRJNApNVEV5T1RBNE5ETXhORm93Z1pjeEN6QUpCZ05WQkFZVEFrTk9NUkF3RGdZRFZRUUlFd2RDWldscWFXNW5NUkF3CkRnWURWUVFIRXdkQ1pXbHFhVzVuTVI4d0hRWURWUVFLRXhaM2VDMXZjbWMwTG1Ob1lXbHViV0ZyWlhJdWIzSm4KTVJJd0VBWURWUVFMRXdsamIyNXpaVzV6ZFhNeEx6QXRCZ05WQkFNVEptTnZibk5sYm5OMWN6RXVjMmxuYmk1MwplQzF2Y21jMExtTm9ZV2x1YldGclpYSXViM0puTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFCmI4aFpkU2hESEsra3ErSjZESnFUVnZucjE3TjVOL2pubTE0dHI4eWFXZllTaDZWd0sxM3JFdW5aRXhrcXZjR3oKVXErNmhWTmhCNXRLZHM3SGZOREpuS05xTUdnd0RnWURWUjBQQVFIL0JBUURBZ2JBTUNrR0ExVWREZ1FpQkNCSgpJMUNERHArclcvUVBVRXlyTDFYemF5SGRYL0tmQ0lKa3hsRTY0bTlxUkRBckJnTlZIU01FSkRBaWdDRHJBaGl3CnVZMjZXREc0c0RqU0R6MDV3bWZ1WElCNFdIZk93VlJvVFB4ZzREQUtCZ2dxaGtqT1BRUURBZ05KQURCR0FpRUEKZ1JLMS9tVENXcTNoZHdpSmpJL09IMVVjT3VsK2lDUTd2NEtwb0pYZnRhZ0NJUUQvVUJkQzI0Qm9xRzZhNGthTwpnY3JpbGhVV00rOTlkRG1IRjlHaUp3YzYxZz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KEkcwRQIgMwyO6gh7Ns1m4cD6weUGBKmTl8UPCy1Jb+353mHDv9cCIQC+g5nk45vNrhzsmnB95BcFMm0n5xbO2LcUlOv3NY+EDwqRCQouUW1US3d5c3BUNTVKY3E4aWVuMUpNeXhja042QW1xZ0JwQ3oydlVaeXA0TlJkWBLeCAgBEi5RbVRLd3lzcFQ1NUpjcThpZW4xSk15eGNrTjZBbXFnQnBDejJ2VVp5cDROUmRYGB4gDCogAg42/wT1z9OBZ956lzM+GMMsz6eWnY7Xe1Td5/SPCl4ygwgKtwcKFnd4LW9yZzIuY2hhaW5tYWtlci5vcmcanActLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ2ZqQ0NBaVNnQXdJQkFnSURBcUJvTUFvR0NDcUdTTTQ5QkFNQ01JR0tNUXN3Q1FZRFZRUUdFd0pEVGpFUQpNQTRHQTFVRUNCTUhRbVZwYW1sdVp6RVFNQTRHQTFVRUJ4TUhRbVZwYW1sdVp6RWZNQjBHQTFVRUNoTVdkM2d0CmIzSm5NaTVqYUdGcGJtMWhhMlZ5TG05eVp6RVNNQkFHQTFVRUN4TUpjbTl2ZEMxalpYSjBNU0l3SUFZRFZRUUQKRXhsallTNTNlQzF2Y21jeUxtTm9ZV2x1YldGclpYSXViM0puTUI0WERUSXpNVEl3TVRBNE5ETXhORm9YRFRJNApNVEV5T1RBNE5ETXhORm93Z1pjeEN6QUpCZ05WQkFZVEFrTk9NUkF3RGdZRFZRUUlFd2RDWldscWFXNW5NUkF3CkRnWURWUVFIRXdkQ1pXbHFhVzVuTVI4d0hRWURWUVFLRXhaM2VDMXZjbWN5TG1Ob1lXbHViV0ZyWlhJdWIzSm4KTVJJd0VBWURWUVFMRXdsamIyNXpaVzV6ZFhNeEx6QXRCZ05WQkFNVEptTnZibk5sYm5OMWN6RXVjMmxuYmk1MwplQzF2Y21jeUxtTm9ZV2x1YldGclpYSXViM0puTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFCitEeXRCTmV2RUo0WFhEL3RacXQyQlk3Kytza0NDajB0Q2owN09JMkZmVWNVSU05Yk1QRnJZaDBFRGdKQlMvS0oKME5uZUVlTzRvK3ZBQndweDdsRENlYU5xTUdnd0RnWURWUjBQQVFIL0JBUURBZ2JBTUNrR0ExVWREZ1FpQkNBRQpnYVZFN09xR0grOGMyR3N6dGxpMXZsZGs2bFhJNDJyMFMwWVBjTU94VlRBckJnTlZIU01FSkRBaWdDQjE2YktGCmtnSmtpb2hzUEo2ZFNjVDRjbjJtUXZUQUhyM2ZaV3FsNXc3azFEQUtCZ2dxaGtqT1BRUURBZ05JQURCRkFpRUEKa3ZvKzI4bEMxR3FnSXJOTmFYeVZnV3FnUkJBQmJ1K2VVd2EwcEhlbE5Hb0NJQlNPbEYvbW9yRFEwSXJ1OGRlZwpCR0cxR0V0bmFseFlFTnZtVzJ6VlA3ZnYKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQoSRzBFAiAyHVnNphEchdBVTDMaa6FcTy9t3M2qnckb9lE4gG60QAIhAOSscalvu3scGol9YMMQGVDbXyoSgHTGZv52d3wVH97XEAM=\"}}},\"rwset_list\":[{\"tx_id\":\"17a05aa42d8d673cca4160edc45863596289046afae0442c9a5bcd6bc5347359\",\"tx_reads\":[{\"key\":\"X19hY2NvdW50X3ByZWZpeF9fMTcxMjYyMzQ3YTU5ZmRlZDkyMDIxYTMyNDIxYTVkYWQwNTQyNGUwMw==\",\"value\":\"OTk5OTk5OTk5OTgxNzMxNg==\",\"contract_name\":\"ACCOUNT_MANAGER\"},{\"key\":\"X19mcm96ZW5fYWNjb3VudF9fMTcxMjYyMzQ3YTU5ZmRlZDkyMDIxYTMyNDIxYTVkYWQwNTQyNGUwMw==\",\"contract_name\":\"ACCOUNT_MANAGER\"},{\"key\":\"MmNhZTNmMmE5OTU4YjEwNzA1NTM2MDI1OTc0NTMxYTJiYjI5NTQwMmRkYjdmYjYwMWVlNjdjNDMxODhjMTM2Mw==\",\"value\":\"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNlRENDQWg2Z0F3SUJBZ0lERFdIZk1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTVM1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3hMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2daRXhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jeExtTm9ZV2x1YldGclpYSXViM0puCk1ROHdEUVlEVlFRTEV3WmpiR2xsYm5ReExEQXFCZ05WQkFNVEkyTnNhV1Z1ZERFdWMybG5iaTUzZUMxdmNtY3gKTG1Ob1lXbHViV0ZyWlhJdWIzSm5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVadTcwb0xRWQp2UEptNnZlUFdMbENWeTJHOHpqYUQxL2tpYUpuMnNyRnc3WVR2cDNYV2d5OVJVM1ZRNnJsa004VExaYWY5Z2NRCmdScWFoSEJnaTZ0T1FLTnFNR2d3RGdZRFZSMFBBUUgvQkFRREFnYkFNQ2tHQTFVZERnUWlCQ0FvNGxwakRiNXAKSmdCc2JBc3U5aXEwQlM1V3p3N0IvMy9kelM0anpadEdTakFyQmdOVkhTTUVKREFpZ0NBUFJxKy8xd1FQajhBawplVkl5bDhENmkwZGdxdnh5NWV1QytERjVXVnVVTnpBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBakdlZ0pndWQ1CnZPU0plVktENzdyUzFwOWE5TytQQU1UM3ptbWd6MlJZWndJaEFPNDE4Z3V2NUlhckFJMmt1MXlGbTVQK2FmYWQKeW1lNnp2c1RVbEdhOHhLZgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==\",\"contract_name\":\"CERT_MANAGE\"},{\"key\":\"Q29udHJhY3Q6RVZNX0VSQzIwXzE=\",\"contract_name\":\"CONTRACT_MANAGE\"}],\"tx_writes\":[{\"key\":\"MDM=\",\"value\":\"AzsuPJ/QgDzoAAAA\",\"contract_name\":\"882e540f48d6aeb9c64f9affc748624a957d3d4f\"},{\"key\":\"X19hY2NvdW50X3ByZWZpeF9fMTcxMjYyMzQ3YTU5ZmRlZDkyMDIxYTMyNDIxYTVkYWQwNTQyNGUwMw==\",\"value\":\"OTk5OTk5OTk5OTkxNzQ2Mw==\",\"contract_name\":\"ACCOUNT_MANAGER\"},{\"key\":\"Q29udHJhY3Q6ODgyZTU0MGY0OGQ2YWViOWM2NGY5YWZmYzc0ODYyNGE5NTdkM2Q0Zg==\",\"value\":\"CgtFVk1fRVJDMjBfMRIDMS4wGAUqqwEKFnd4LW9yZzEuY2hhaW5tYWtlci5vcmcQARogLK4/KplYsQcFU2All0UxorspVALdt/tgHuZ8QxiME2MiI2NsaWVudDEuc2lnbi53eC1vcmcxLmNoYWlubWFrZXIub3JnKgZDTElFTlQyQDI4ZTI1YTYzMGRiZTY5MjYwMDZjNmMwYjJlZjYyYWI0MDUyZTU2Y2YwZWMxZmY3ZmRkY2QyZTIzY2Q5YjQ2NGEyKDg4MmU1NDBmNDhkNmFlYjljNjRmOWFmZmM3NDg2MjRhOTU3ZDNkNGY=\",\"contract_name\":\"CONTRACT_MANAGE\"},{\"key\":\"Q29udHJhY3Q6RVZNX0VSQzIwXzE=\",\"value\":\"CgtFVk1fRVJDMjBfMRIDMS4wGAUqqwEKFnd4LW9yZzEuY2hhaW5tYWtlci5vcmcQARogLK4/KplYsQcFU2All0UxorspVALdt/tgHuZ8QxiME2MiI2NsaWVudDEuc2lnbi53eC1vcmcxLmNoYWlubWFrZXIub3JnKgZDTElFTlQyQDI4ZTI1YTYzMGRiZTY5MjYwMDZjNmMwYjJlZjYyYWI0MDUyZTU2Y2YwZWMxZmY3ZmRkY2QyZTIzY2Q5YjQ2NGEyKDg4MmU1NDBmNDhkNmFlYjljNjRmOWFmZmM3NDg2MjRhOTU3ZDNkNGY=\",\"contract_name\":\"CONTRACT_MANAGE\"},{\"key\":\"Q29udHJhY3RCeXRlQ29kZTpFVk1fRVJDMjBfMQ==\",\"value\":\"YGBgQFJgBDYQYQCZV2AANXwBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJAEY/////8WgGMG/d4DFGEAnleAYwlep7MUYQEsV4BjGBYN3RRhAYZXgGMjuHLdFGEBr1eAYzE85WcUYQIoV4BjcKCCMRRhAldXgGOV2JtBFGECpFeAY6kFnLsUYQMyV4Bj3WLtPhRhA4xXW2AAgP1bNBVhAKlXYACA/VthALFhA/hWW2BAUYCAYCABgoEDglKDgYFRgVJgIAGRUIBRkGAgAZCAg4NgAFuDgRAVYQDxV4CCAVGBhAFSYCCBAZBQYQDWVltQUFBQkFCQgQGQYB8WgBVhAR5XgIIDgFFgAYNgIANhAQAKAxkWgVJgIAGRUFtQklBQUGBAUYCRA5DzWzQVYQE3V2AAgP1bYQFsYASAgDVz//////////////////////////8WkGAgAZCRkIA1kGAgAZCRkFBQYQQxVltgQFGAghUVFRWBUmAgAZFQUGBAUYCRA5DzWzQVYQGRV2AAgP1bYQGZYQUjVltgQFGAgoFSYCABkVBQYEBRgJEDkPNbNBVhAbpXYACA/VthAg5gBICANXP//////////////////////////xaQYCABkJGQgDVz//////////////////////////8WkGAgAZCRkIA1kGAgAZCRkFBQYQUpVltgQFGAghUVFRWBUmAgAZFQUGBAUYCRA5DzWzQVYQIzV2AAgP1bYQI7YQhOVltgQFGAgmD/FmD/FoFSYCABkVBQYEBRgJEDkPNbNBVhAmJXYACA/VthAo5gBICANXP//////////////////////////xaQYCABkJGQUFBhCFNWW2BAUYCCgVJgIAGRUFBgQFGAkQOQ81s0FWECr1dgAID9W2ECt2EIm1ZbYEBRgIBgIAGCgQOCUoOBgVGBUmAgAZFQgFGQYCABkICDg2AAW4OBEBVhAvdXgIIBUYGEAVJgIIEBkFBhAtxWW1BQUFCQUJCBAZBgHxaAFWEDJFeAggOAUWABg2AgA2EBAAoDGRaBUmAgAZFQW1CSUFBQYEBRgJEDkPNbNBVhAz1XYACA/VthA3JgBICANXP//////////////////////////xaQYCABkJGQgDWQYCABkJGQUFBhCNRWW2BAUYCCFRUVFYFSYCABkVBQYEBRgJEDkPNbNBVhA5dXYACA/VthA+JgBICANXP//////////////////////////xaQYCABkJGQgDVz//////////////////////////8WkGAgAZCRkFBQYQq4VltgQFGAgoFSYCABkVBQYEBRgJEDkPNbYECAUZCBAWBAUoBgE4FSYCABf1Rlc3RDaGFpbk1ha2VyVG9rZW4AAAAAAAAAAAAAAAAAgVJQgVZbYACBYAFgADNz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYACFc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIIGQVVCCc///////////////////////////FjNz//////////////////////////8Wf4xb4eXr7H1b0U9xQn0ehPPdAxTA97IpHlsgCsjHw7klhGBAUYCCgVJgIAGRUFBgQFGAkQOQo2ABkFCSkVBQVltgA1SBVltgAIBgAWAAhnP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBgADNz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVJBQgmAAgIdz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVBAVgBVhBflXUIKBEBVbgBVhBoNXUGAAgIVz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVINgAICHc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIFQBEBVbFWEIQVeCYACAhnP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBgAIKCVAGSUFCBkFVQgmAAgIdz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYACCglQDklBQgZBVUH///////////////////////////////////////////4EQFWEH01eCYAFgAIdz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYAAzc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIGAAgoJUA5JQUIGQVVBbg3P//////////////////////////xaFc///////////////////////////Fn/d8lKtG+LIm2nCsGj8N42qlSun8WPEoRYo9VpN9SOz74VgQFGAgoFSYCABkVBQYEBRgJEDkKNgAZFQYQhGVltgAJFQW1CTklBQUFZbYBKBVltgAIBgAINz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVJBQkZBQVltgQIBRkIEBYEBSgGADgVJgIAF/VENNAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACBUlCBVltgAIFgAIAzc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIFQQFYAVYQmiV1BgAICEc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIFSCYACAhnP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBUARAVWxVhCq1XgWAAgDNz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgYACCglQDklBQgZBVUIFgAICFc///////////////////////////FnP//////////////////////////xaBUmAgAZCBUmAgAWAAIGAAgoJUAZJQUIGQVVCCc///////////////////////////FjNz//////////////////////////8Wf93yUq0b4sibacKwaPw3jaqVK6fxY8ShFij1Wk31I7PvhGBAUYCCgVJgIAGRUFBgQFGAkQOQo2ABkFBhCrJWW2AAkFBbkpFQUFZbYABgAWAAhHP//////////////////////////xZz//////////////////////////8WgVJgIAGQgVJgIAFgACBgAINz//////////////////////////8Wc///////////////////////////FoFSYCABkIFSYCABYAAgVJBQkpFQUFYAoWVienpyMFggJZE/R6lfV/obNOchph0QuRIBwjsDd47+Oai3P2pV+ccAKQ==\",\"contract_name\":\"CONTRACT_MANAGE\"},{\"key\":\"MDM=\",\"value\":\"AzsuPJ/QgDzoAAAA\",\"contract_name\":\"EVM_ERC20_1\"},{\"key\":\"ZTIzODIxYWQ2Mjc5NDliMjczMTYxMjEyNTMzZGE5MWYxNGZhOWZlZDVjYThjOWFiNWQzNjNjMzZiYWE2YzI0ZA==\",\"value\":\"AzsuPJ/QgDzoAAAA\",\"contract_name\":\"EVM_ERC20_1\"}]}]}"
var BlockResultJson = "{\"preBlockHash\":\"a5dc0b4a974cf055ef03b3c32570b8a2dad101dc22df27e11cc07d81dd1a301d\",\"blockHash\":\"020e36ff04f5cfd38167de7a97333e18c32ccfa7969d8ed77b54dde7f48f0a5e\",\"blockHeight\":30,\"blockVersion\":2030200,\"orgId\":\"wx-org3.chainmaker.org\",\"timestamp\":1702460320,\"blockDag\":\"{\\\"vertexes\\\":[{}]}\",\"dagHash\":\"08da7c45cb204377e7e42249cda5713fa865116ddbb4cb5a1949b2e5b438a6ab\",\"txCount\":1,\"signature\":\"MEQCIFBzMeUG/A67UaVqQUParBgMBimG+bKd9AtesPRnqlWuAiBkTrgev/YdVlmUPg20r4RGdDpzuafib+IOvUTW93zr/w==\",\"rwSetHash\":\"7fd7626c09e09e7863385f665b0a325bf1d684553721efbf520d30ca687793ca\",\"txRootHash\":\"40a5e4734a4fce19530e9673e3c3470ba1a274611a574aeb64fd55de07edb789\",\"proposerId\":\"consensus1.sign.wx-org3.chainmaker.org\",\"proposerAddr\":\"442dbafee70852e741e01b4b98dd2cf90ece1d96\",\"consensusArgs\":\"\",\"delayUpdateStatus\":0,\"createdAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}"

func TestMain(m *testing.M) {
	// 初始化数据库配置
	redisCfg, err := db.InitRedisContainer()
	if err != nil {
		return
	}
	_, err = db.InitMySQLContainer()
	if err != nil {
		return
	}

	cache.InitRedis(redisCfg)
	// 运行其他测试
	os.Exit(m.Run())
}

func TestDealBlockInfo(t *testing.T) {
	blockInfo := &common.BlockInfo{}
	err := json.Unmarshal([]byte(BlockJson), blockInfo)
	if err != nil {
		return
	}
	blockDB := &db.Block{}
	err = json.Unmarshal([]byte(BlockResultJson), blockDB)
	if err != nil {
		return
	}

	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	chainId := blockInfo.Block.Header.ChainId
	err = SetMemberInfoCache(chainId, "SHA256", txInfo.Sender.Signer)
	if err != nil {
		return
	}

	type args struct {
		blockInfo *common.BlockInfo
		hashType  string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Block
		wantErr bool
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				blockInfo: blockInfo,
				hashType:  "SHA256",
			},
			want:    blockDB,
			wantErr: false,
		},
		{
			name: "Test case 2: Invalid blockInfo",
			args: args{
				blockInfo: nil,
				hashType:  "SHA256",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DealBlockInfo(tt.args.blockInfo, tt.args.hashType)
			if (err != nil) != tt.wantErr {
				t.Errorf("DealBlockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Block{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("DealBlockInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSetBlockListCache(t *testing.T) {
	blockDB := &db.Block{}
	err := json.Unmarshal([]byte(BlockResultJson), blockDB)
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		blockList []*db.Block
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				chainId:   "testchain1",
				blockList: []*db.Block{blockDB},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbhandle.SetLatestBlockListCache(tt.args.chainId, tt.args.blockList)
		})
	}
}

func TestGetBlockListFromRedis(t *testing.T) {
	TestSetBlockListCache(t)
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Block
		wantErr bool
	}{
		{
			name: "Test case 1: Valid chainId and blockList",
			args: args{
				chainId: "testchain1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbhandle.GetLatestBlockListCache(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockListFromRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GetBlockListFromRedis() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealBlockInfo1(t *testing.T) {
	blockInfo := &common.BlockInfo{}
	err := json.Unmarshal([]byte(BlockJson), blockInfo)
	if err != nil {
		return
	}
	blockDB := &db.Block{}
	err = json.Unmarshal([]byte(BlockResultJson), blockDB)
	if err != nil {
		return
	}

	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	chainId := blockInfo.Block.Header.ChainId
	err = SetMemberInfoCache(chainId, "SHA256", txInfo.Sender.Signer)
	if err != nil {
		return
	}

	type args struct {
		blockInfo *common.BlockInfo
		hashType  string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Block
		wantErr bool
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				blockInfo: blockInfo,
				hashType:  "SHA256",
			},
			want:    blockDB,
			wantErr: false,
		},
		{
			name: "Test case 2: Invalid blockInfo",
			args: args{
				blockInfo: nil,
				hashType:  "SHA256",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DealBlockInfo(tt.args.blockInfo, tt.args.hashType)
			if (err != nil) != tt.wantErr {
				t.Errorf("DealBlockInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Block{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("DealBlockInfo() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestBuildLatestBlockListCache(t *testing.T) {
	type args struct {
		chainId  string
		modBlock *db.Block
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				chainId: ChainId,
				modBlock: &db.Block{
					BlockHeight:  12,
					PreBlockHash: "12",
					BlockHash:    "12",
					OrgId:        "12",
					Timestamp:    1223323,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildLatestBlockListCache(tt.args.chainId, tt.args.modBlock)
		})
	}
}
