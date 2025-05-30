package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// PluginMetadata represents plugin metadata
type PluginMetadata struct {
    ID           string            `json:"id"`
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description"`
    Author       string            `json:"author"`
    Endpoints    []EndpointConfig  `json:"endpoints"`
    Dependencies []string          `json:"dependencies"`
    Permissions  []string          `json:"permissions"`
    Compatibility map[string]string `json:"compatibility"`
}

// EndpointConfig represents an API endpoint configuration
type EndpointConfig struct {
    Path        string `json:"path"`
    Method      string `json:"method"`
    Handler     string `json:"handler"`
    Description string `json:"description"`
}

// ClusterInfo represents cluster information
type ClusterInfo struct {
    ClusterName string    `json:"clusterName"`
    Status      string    `json:"status"`
    Message     string    `json:"message"`
    LastUpdated time.Time `json:"lastUpdated"`
    NodeCount   int       `json:"nodeCount"`
    Namespace   string    `json:"namespace,omitempty"`
    Region      string    `json:"region,omitempty"`
}

// OnboardRequest represents cluster onboarding request
type OnboardRequest struct {
    ClusterName string `json:"clusterName" binding:"required"`
    Namespace   string `json:"namespace"`
    Region      string `json:"region"`
    AutoSetup   bool   `json:"autoSetup"`
}

// DetachRequest represents cluster detach request
type DetachRequest struct {
    ClusterName string `json:"clusterName" binding:"required"`
    ForceDetach bool   `json:"forceDetach"`
    CleanupData bool   `json:"cleanupData"`
}

// TestClusterPlugin represents the demo plugin instance
type TestClusterPlugin struct {
    initialized bool
    clusters    map[string]*ClusterInfo
    startTime   time.Time
}

// Initialize initializes the plugin
func (p *TestClusterPlugin) Initialize(config map[string]interface{}) error {
    if p.initialized {
        fmt.Println("⚠️ Plugin already initialized")
        return nil
    }

    p.clusters = make(map[string]*ClusterInfo)
    p.startTime = time.Now()
    
    // Add some demo clusters
    p.clusters["prod-cluster-east"] = &ClusterInfo{
        ClusterName: "prod-cluster-east",
        Status:      "ready",
        Message:     "Production cluster - East region",
        LastUpdated: time.Now().Add(-10 * time.Minute),
        NodeCount:   5,
        Namespace:   "kubestellar-system",
        Region:      "us-east-1",
    }
    
    p.clusters["staging-cluster"] = &ClusterInfo{
        ClusterName: "staging-cluster",
        Status:      "ready",
        Message:     "Staging environment ready",
        LastUpdated: time.Now().Add(-5 * time.Minute),
        NodeCount:   2,
        Namespace:   "kubestellar-system",
        Region:      "us-west-2",
    }

    p.clusters["dev-cluster-1"] = &ClusterInfo{
        ClusterName: "dev-cluster-1",
        Status:      "pending",
        Message:     "Development cluster initializing",
        LastUpdated: time.Now().Add(-2 * time.Minute),
        NodeCount:   1,
        Namespace:   "kubestellar-system",
        Region:      "eu-west-1",
    }

    p.initialized = true
    fmt.Println("✅ TestClusterPlugin initialized successfully")
    return nil
}

// GetMetadata returns plugin metadata
func (p *TestClusterPlugin) GetMetadata() PluginMetadata {
    return PluginMetadata{
        ID:          "kubestellar-demo-plugin",
        Name:        "KubeStellar Demo Plugin v1.0 - Working Build", // ← This will show successful build
        Version:     "1.0.0",
        Description: "Demonstration plugin for KubeStellar cluster management",
        Author:      "CNCF LFX Mentee",
        Endpoints: []EndpointConfig{
            {Path: "/status", Method: "GET", Handler: "GetClusterStatusHandler", Description: "Get plugin status"},
            {Path: "/clusters", Method: "GET", Handler: "ListClustersHandler", Description: "List clusters"},
            {Path: "/onboard", Method: "POST", Handler: "OnboardClusterHandler", Description: "Onboard cluster"},
            {Path: "/detach", Method: "POST", Handler: "DetachClusterHandler", Description: "Detach cluster"},
            {Path: "/health", Method: "GET", Handler: "HealthCheckHandler", Description: "Health check"},
        },
        Dependencies: []string{"kubectl", "clusteradm"},
        Permissions:  []string{"cluster.read", "cluster.write", "namespace.create"},
        Compatibility: map[string]string{
            "kubestellar": ">=0.21.0",
            "go":          ">=1.21",
            "kubernetes":  ">=1.24",
        },
    }
}

// GetHandlers returns the plugin's HTTP handlers
func (p *TestClusterPlugin) GetHandlers() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc{
        "GetClusterStatusHandler": p.GetClusterStatusHandler,
        "ListClustersHandler":     p.ListClustersHandler,
        "OnboardClusterHandler":   p.OnboardClusterHandler,
        "DetachClusterHandler":    p.DetachClusterHandler,
        "HealthCheckHandler":      p.HealthCheckHandler,
    }
}

// GetClusterStatusHandler returns plugin and cluster status
func (p *TestClusterPlugin) GetClusterStatusHandler(c *gin.Context) {
    fmt.Println("📊 GetClusterStatusHandler called")
    
    clusters := make([]*ClusterInfo, 0, len(p.clusters))
    ready, pending, failed := 0, 0, 0
    
    for _, cluster := range p.clusters {
        clusters = append(clusters, cluster)
        switch cluster.Status {
        case "ready":
            ready++
        case "pending":
            pending++
        case "failed":
            failed++
        }
    }

    metadata := p.GetMetadata()
    
    response := gin.H{
        "plugin":    metadata.Name,
        "version":   metadata.Version,
        "author":    metadata.Author,
        "uptime":    time.Since(p.startTime).String(),
        "timestamp": time.Now(),
        "clusters":  clusters,
        "summary": gin.H{
            "total":   len(clusters),
            "ready":   ready,
            "pending": pending,
            "failed":  failed,
        },
        "metadata": metadata,
    }

    fmt.Printf("✅ Returning cluster status: %d clusters\n", len(clusters))
    c.JSON(http.StatusOK, response)
}

// ListClustersHandler returns list of all clusters
func (p *TestClusterPlugin) ListClustersHandler(c *gin.Context) {
    fmt.Println("📋 ListClustersHandler called")
    
    clusters := make([]*ClusterInfo, 0, len(p.clusters))
    for _, cluster := range p.clusters {
        clusters = append(clusters, cluster)
    }

    c.JSON(http.StatusOK, gin.H{
        "clusters":  clusters,
        "count":     len(clusters),
        "timestamp": time.Now(),
    })
}

// OnboardClusterHandler handles cluster onboarding
func (p *TestClusterPlugin) OnboardClusterHandler(c *gin.Context) {
    fmt.Println("🚀 OnboardClusterHandler called")
    
    var req OnboardRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request",
            "details": err.Error(),
        })
        return
    }

    // Check if cluster already exists
    if _, exists := p.clusters[req.ClusterName]; exists {
        c.JSON(http.StatusConflict, gin.H{
            "error":   "Cluster already exists",
            "cluster": req.ClusterName,
        })
        return
    }

    // Set defaults
    namespace := req.Namespace
    if namespace == "" {
        namespace = "kubestellar-system"
    }
    
    region := req.Region
    if region == "" {
        regions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-south-1"}
        region = regions[rand.Intn(len(regions))]
    }

    // Simulate onboarding process
    status := "pending"
    message := fmt.Sprintf("Cluster onboarding initiated for region %s", region)
    nodeCount := rand.Intn(5) + 1

    // 30% chance of immediate success for demo
    if rand.Float32() < 0.3 {
        status = "ready"
        message = fmt.Sprintf("Cluster onboarded successfully in %s", region)
    }

    // Add to clusters
    p.clusters[req.ClusterName] = &ClusterInfo{
        ClusterName: req.ClusterName,
        Status:      status,
        Message:     message,
        LastUpdated: time.Now(),
        NodeCount:   nodeCount,
        Namespace:   namespace,
        Region:      region,
    }

    fmt.Printf("✅ Cluster '%s' onboarding started\n", req.ClusterName)
    
    c.JSON(http.StatusOK, gin.H{
        "message":     fmt.Sprintf("Cluster '%s' onboarding started", req.ClusterName),
        "clusterName": req.ClusterName,
        "status":      status,
        "namespace":   namespace,
        "region":      region,
        "nodeCount":   nodeCount,
        "timestamp":   time.Now(),
    })
}

// DetachClusterHandler handles cluster detachment
func (p *TestClusterPlugin) DetachClusterHandler(c *gin.Context) {
    fmt.Println("🗑️ DetachClusterHandler called")
    
    var req DetachRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request",
            "details": err.Error(),
        })
        return
    }

    // Check if cluster exists
    cluster, exists := p.clusters[req.ClusterName]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "Cluster not found",
            "cluster": req.ClusterName,
        })
        return
    }

    // Remove cluster
    delete(p.clusters, req.ClusterName)

    fmt.Printf("✅ Cluster '%s' detached successfully\n", req.ClusterName)
    
    c.JSON(http.StatusOK, gin.H{
        "message":        fmt.Sprintf("Cluster '%s' detached successfully", req.ClusterName),
        "clusterName":    req.ClusterName,
        "previousStatus": cluster.Status,
        "cleanupData":    req.CleanupData,
        "timestamp":      time.Now(),
    })
}

// HealthCheckHandler returns plugin health status
func (p *TestClusterPlugin) HealthCheckHandler(c *gin.Context) {
    uptime := time.Since(p.startTime)
    healthy := uptime > 0 && p.initialized
    
    status := "healthy"
    if !healthy {
        status = "unhealthy"
    }

    readyClusters := 0
    for _, cluster := range p.clusters {
        if cluster.Status == "ready" {
            readyClusters++
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    status,
        "uptime":    uptime.String(),
        "clusters":  len(p.clusters),
        "ready":     readyClusters,
        "timestamp": time.Now(),
        "checks": gin.H{
            "initialized":     p.initialized,
            "uptime_ok":       uptime > 0,
            "clusters_loaded": len(p.clusters) > 0,
            "ready_clusters":  readyClusters,
        },
    })
}

// Health implements the health check interface
func (p *TestClusterPlugin) Health() error {
    if !p.initialized {
        return fmt.Errorf("plugin not initialized")
    }
    return nil
}

// Cleanup cleans up plugin resources
func (p *TestClusterPlugin) Cleanup() error {
    fmt.Println("🧹 TestClusterPlugin cleanup initiated")
    p.clusters = make(map[string]*ClusterInfo)
    p.initialized = false
    fmt.Println("✅ TestClusterPlugin cleanup completed")
    return nil
}

// NewPlugin creates a new plugin instance (required export)
func NewPlugin() interface{} {
    fmt.Println("🏗️ Creating new TestClusterPlugin instance")
    plugin := &TestClusterPlugin{}
    return plugin
}
