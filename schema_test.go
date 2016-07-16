package sqlstmt

import (
	"testing"
)

func TestSchemaDefaults(t *testing.T) {
	dialect1 := NewDialect("postgres")
	dialect2 := NewDialect("ql")
	convention1 := ConventionSame
	convention2 := ConventionSnake

	tests := []struct {
		defaultConvention  Convention
		defaultDialect     Dialect
		schemaConvention   Convention
		schemaDialect      Dialect
		expectedConvention Convention
		expectedDialect    Dialect
	}{
		{
			defaultConvention:  convention1,
			defaultDialect:     dialect1,
			expectedConvention: convention1,
			expectedDialect:    dialect1,
		},
		{
			defaultConvention:  convention2,
			defaultDialect:     dialect1,
			schemaConvention:   convention1,
			schemaDialect:      dialect2,
			expectedConvention: convention1,
			expectedDialect:    dialect2,
		},
	}

	resetDefaults := func() {
		DefaultSchema = &Schema{}
	}
	defer resetDefaults()

	for _, tt := range tests {
		resetDefaults()
		DefaultSchema.Convention = tt.defaultConvention
		DefaultSchema.Dialect = tt.defaultDialect
		schema := &Schema{
			Convention: tt.schemaConvention,
			Dialect:    tt.schemaDialect,
		}
		if schema.convention().ColumnName("XyzAbc") != tt.expectedConvention.ColumnName("XyzAbc") {
			t.Errorf("unexpected convention: %v, %v", schema.convention(), tt.expectedConvention)
		}
		if schema.dialect().Name() != tt.expectedDialect.Name() {
			t.Error("unexpected dialect")
		}
	}
}
