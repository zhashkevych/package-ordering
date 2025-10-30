package api

import (
	"net/http"
	"sort"

	"github.com/zhashkevych/package-ordering/internal/order"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router    *gin.Engine
	packSizes []int
}

func NewServer(initialPacks []int) *Server {
	s := &Server{
		router:    gin.Default(),
		packSizes: dedupeAndSortDesc(initialPacks),
	}
	s.registerRoutes()
	return s
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Run(addr ...string) error {
	return s.router.Run(addr...)
}

func (s *Server) registerRoutes() {
	r := s.router

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/packs", func(c *gin.Context) {
		cp := append([]int{}, s.packSizes...)
		sort.Ints(cp)
		c.JSON(http.StatusOK, gin.H{"packSizes": cp})
	})

	r.PUT("/packs", func(c *gin.Context) {
		var req SetPacksRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		if len(req.PackSizes) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "packSizes cannot be empty"})
			return
		}
		for _, v := range req.PackSizes {
			if v <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "pack sizes must be positive"})
				return
			}
		}
		s.packSizes = dedupeAndSortDesc(req.PackSizes)
		c.JSON(http.StatusOK, gin.H{"packSizes": s.packSizes})
	})

	r.POST("/calculate", func(c *gin.Context) {
		var req CalculateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		if req.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be > 0"})
			return
		}

		allocation := order.CalculatePacks(req.Amount, s.packSizes)
		if allocation == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no feasible allocation"})
			return
		}

		totalItems := 0
		totalPacks := 0
		for sz, qty := range allocation {
			totalItems += sz * qty
			totalPacks += qty
		}
		resp := CalculateResponse{
			Amount:     req.Amount,
			PackSizes:  append([]int{}, s.packSizes...),
			Allocation: allocation,
			TotalItems: totalItems,
			TotalPacks: totalPacks,
			Overfill:   totalItems - req.Amount,
		}
		c.JSON(http.StatusOK, resp)
	})

	// Serve static UI from ./web without conflicting with API routes
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/ui/") })
	r.Static("/ui", "./web")
}

func dedupeAndSortDesc(in []int) []int {
	m := make(map[int]struct{}, len(in))
	for _, v := range in {
		if v > 0 {
			m[v] = struct{}{}
		}
	}
	out := make([]int, 0, len(m))
	for v := range m {
		out = append(out, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(out)))
	return out
}
