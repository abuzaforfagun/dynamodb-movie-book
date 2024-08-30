package models

type TopRatedMovies struct {
	PK     string
	SK     string
	GSI_PK string
	GSI_SK string
	Movies []*MovieShortInformation
}

func NewTopRatedMovies(movies []*MovieShortInformation, numberOfMovies int) *TopRatedMovies {
	uniqueMovies := UniqueMovies(movies)
	if len(uniqueMovies) > numberOfMovies {
		uniqueMovies = (uniqueMovies)[:numberOfMovies]
	}
	return &TopRatedMovies{
		PK:     "TOP-RATED-MOVIE",
		SK:     "TOP-RATED-MOVIE",
		GSI_PK: "TOP-RATED-MOVIE",
		GSI_SK: "TOP-RATED-MOVIE",
		Movies: uniqueMovies,
	}
}

func UniqueMovies(movies []*MovieShortInformation) []*MovieShortInformation {
	seen := make(map[string]bool)
	var uniqueList []*MovieShortInformation

	for _, movie := range movies {
		if _, exists := seen[movie.Id]; !exists {
			seen[movie.Id] = true
			uniqueList = append(uniqueList, movie)
		}
	}

	return uniqueList
}

type SortByScore []*MovieShortInformation

func (a SortByScore) Len() int           { return len(a) }
func (a SortByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
