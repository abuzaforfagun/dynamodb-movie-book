package movies_handler

import (
	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
}

func New() *MoviesHandler {
	return &MoviesHandler{}
}

func (mh *MoviesHandler) GetAllMovies(c *gin.Context)    {}
func (mh *MoviesHandler) GetMovieDetails(c *gin.Context) {}
func (mh *MoviesHandler) SearchMovies(c *gin.Context)    {}
func (mh *MoviesHandler) GetMovieByGenre(c *gin.Context) {}

func (mh *MoviesHandler) AddMovie(c *gin.Context)    {}
func (mh *MoviesHandler) DeleteMovie(c *gin.Context) {}
