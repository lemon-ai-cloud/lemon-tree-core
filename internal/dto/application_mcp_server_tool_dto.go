// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，确保数据格式的一致性
package dto

// ApplicationMcpServerToolDto ApplicationMCP工具数据传输对象
// 用于在业务逻辑层和HTTP处理层之间传递数据
type ApplicationMcpServerToolDto struct {
	ID                           string `json:"id"`                               // 主键ID
	ApplicationID                string `json:"application_id"`                   // 所属应用ID
	ApplicationMcpServerConfigID string `json:"application_mcp_server_config_id"` // 所属MCP服务配置ID
	Name                         string `json:"name"`                             // 名称
	Title                        string `json:"title"`                            // 工具标题
	Description                  string `json:"description"`                      // 描述
	CreatedAt                    string `json:"created_at"`                       // 创建时间
	UpdatedAt                    string `json:"updated_at"`                       // 更新时间
}
