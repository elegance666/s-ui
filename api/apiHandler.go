package api

import (
	"s-ui/util/common"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	ApiService
	apiv2 *APIv2Handler
}

func NewAPIHandler(g *gin.RouterGroup, a2 *APIv2Handler) {
	a := &APIHandler{
		apiv2: a2,
	}
	a.initRouter(g)
}

func (a *APIHandler) initRouter(g *gin.RouterGroup) {
	g.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if !strings.HasSuffix(path, "login") && !strings.HasSuffix(path, "logout") {
			checkLogin(c)
		}
	})
	g.POST("/:postAction", a.postHandler)
	g.GET("/:getAction", a.getHandler)
	g.POST("/addClient", a.addClientHandler)
    g.POST("/updateClient", a.updateClientHandler)
    g.POST("/delClient", a.delClientHandler)
}

func (a *APIHandler) postHandler(c *gin.Context) {
	loginUser := GetLoginUser(c)
	action := c.Param("postAction")

	switch action {
	case "login":
		a.ApiService.Login(c)
	case "changePass":
		a.ApiService.ChangePass(c)
	case "save":
		a.ApiService.Save(c, loginUser)
	case "restartApp":
		a.ApiService.RestartApp(c)
	case "restartSb":
		a.ApiService.RestartSb(c)
	case "linkConvert":
		a.ApiService.LinkConvert(c)
	case "importdb":
		a.ApiService.ImportDb(c)
	case "addToken":
		a.ApiService.AddToken(c)
		a.apiv2.ReloadTokens()
	case "deleteToken":
		a.ApiService.DeleteToken(c)
		a.apiv2.ReloadTokens()
	default:
		jsonMsg(c, "failed", common.NewError("unknown action: ", action))
	}
}

func (a *APIHandler) getHandler(c *gin.Context) {
	action := c.Param("getAction")

	switch action {
	case "logout":
		a.ApiService.Logout(c)
	case "load":
		a.ApiService.LoadData(c)
	case "inbounds", "outbounds", "endpoints", "tls", "clients", "config":
		err := a.ApiService.LoadPartialData(c, []string{action})
		if err != nil {
			jsonMsg(c, action, err)
		}
		return
	case "users":
		a.ApiService.GetUsers(c)
	case "settings":
		a.ApiService.GetSettings(c)
	case "stats":
		a.ApiService.GetStats(c)
	case "status":
		a.ApiService.GetStatus(c)
	case "onlines":
		a.ApiService.GetOnlines(c)
	case "logs":
		a.ApiService.GetLogs(c)
	case "changes":
		a.ApiService.CheckChanges(c)
	case "keypairs":
		a.ApiService.GetKeypairs(c)
	case "getdb":
		a.ApiService.GetDb(c)
	case "tokens":
		a.ApiService.GetTokens(c)
	default:
		jsonMsg(c, "failed", common.NewError("unknown action: ", action))
	}
}

func (a *APIHandler) addClientHandler(c *gin.Context) {
    var request AddClientRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    inboundID := request.ID
    client := request.Settings.Clients[0]

    err := addClientToInbound(inboundID, client)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to add client"})
        return
    }

    c.JSON(200, gin.H{"message": "Client added successfully"})
}

func (a *APIHandler) updateClientHandler(c *gin.Context) {
    var request UpdateClientRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    inboundID := request.ID
    client := request.Settings.Clients[0]

    err := updateClientInInbound(inboundID, client)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to update client"})
        return
    }

    c.JSON(200, gin.H{"message": "Client updated successfully"})
}

func (a *APIHandler) delClientHandler(c *gin.Context) {
    var request DeleteClientRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    inboundID := request.ID
    clientID := request.ClientID

    err := deleteClientFromInbound(inboundID, clientID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to delete client"})
        return
    }

    c.JSON(200, gin.H{"message": "Client deleted successfully"})
}
