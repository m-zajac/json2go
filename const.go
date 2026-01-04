package json2go

const (
	// General options
	DefaultRootName = "Document"

	// Extraction heuristics
	DefaultSimilarityThreshold  = 0.7
	DefaultMinSubsetSize        = 2
	DefaultMinSubsetOccurrences = 2
	DefaultMinAddedFields       = 2

	// Map conversion
	DefaultMakeMaps                  = true
	DefaultMakeMapsWhenMinAttributes = 5

	// Feature toggles
	DefaultExtractAllTypes              = true
	DefaultExtractCommonTypes           = true
	DefaultStringPointersWhenKeyMissing = true
	DefaultSkipEmptyKeys                = true
	DefaultTimeAsStr                    = false
)
