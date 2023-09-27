package osuBot

type MapResult struct {
	Id           int     `bson:"id"`
	Url          string  `bson:"url"`
	Version      string  `bson:"version"`
	Title        string  `bson:"title"`
	PP100Percent float64 `bson:"PP_100_percent"`
	PP99Percent  float64 `bson:"PP_99_percent"`
	PP98Percent  float64 `bson:"PP_98_percent"`
	PP97Percent  float64 `bson:"PP_97_percent"`
	PP96Percent  float64 `bson:"PP_96_percent"`
	PP95Percent  float64 `bson:"PP_95_percent"`
	Results      int     `bson:"results"`
}
