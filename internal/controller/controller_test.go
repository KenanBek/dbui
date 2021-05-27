package controller

import (
	"reflect"
	"testing"

	"github.com/kenanbek/dbui/internal"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/suite"
)

// Configure Suite

type ControllerTestSuite struct {
	suite.Suite

	MockCtrl             *gomock.Controller
	EmptyAppConfig       *MockAppConfig
	TwoConnAppConfig     *MockAppConfig
	UnsupportedAppConfig *MockAppConfig
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.MockCtrl = gomock.NewController(suite.T())

	// empty app config
	suite.EmptyAppConfig = NewMockAppConfig(suite.MockCtrl)
	suite.EmptyAppConfig.EXPECT().DataSourceConfigs().Return(map[string]internal.DataSourceConfig{}).AnyTimes()

	// app config with an unsupported type
	dscUnsupported := NewMockDataSourceConfig(suite.MockCtrl)
	dscUnsupported.EXPECT().Type().Return("mycustomsql").AnyTimes()
	dscUnsupported.EXPECT().Alias().Return("conn1").AnyTimes()
	dscUnsupported.EXPECT().DSN().Return("conn1_dsn").AnyTimes()

	suite.UnsupportedAppConfig = NewMockAppConfig(suite.MockCtrl)
	suite.UnsupportedAppConfig.EXPECT().DataSourceConfigs().Return(map[string]internal.DataSourceConfig{
		dscUnsupported.Alias(): dscUnsupported,
	}).AnyTimes()
	suite.UnsupportedAppConfig.EXPECT().Default().Return(dscUnsupported.Alias()).AnyTimes()

	// two connection app config
	dsc1 := NewMockDataSourceConfig(suite.MockCtrl)
	dsc1.EXPECT().Type().Return("mysql").AnyTimes()
	dsc1.EXPECT().Alias().Return("conn1").AnyTimes()
	dsc1.EXPECT().DSN().Return("conn1_dsn").AnyTimes()

	dsc2 := NewMockDataSourceConfig(suite.MockCtrl)
	dsc2.EXPECT().Type().Return("postgresql").AnyTimes()
	dsc2.EXPECT().Alias().Return("conn2").AnyTimes()
	dsc2.EXPECT().DSN().Return("conn2_dsn").AnyTimes()

	suite.TwoConnAppConfig = NewMockAppConfig(suite.MockCtrl)
	suite.TwoConnAppConfig.EXPECT().DataSourceConfigs().Return(map[string]internal.DataSourceConfig{
		dsc1.Alias(): dsc1,
		dsc2.Alias(): dsc2,
	}).AnyTimes()
	suite.TwoConnAppConfig.EXPECT().Default().Return(dsc1.Alias()).AnyTimes()
}

func (suite *ControllerTestSuite) TearDownAllSuite() {
	suite.MockCtrl.Finish()
}

func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

// Suite Tests

func (suite *ControllerTestSuite) TestNew() {
	type args struct {
		appConfig internal.AppConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{"nil app config", args{}, ErrEmptyConnection},
		{"empty app config", args{suite.EmptyAppConfig}, ErrEmptyConnection},
		{"unsupported db type", args{suite.UnsupportedAppConfig}, ErrUnsupportedDatabaseType},
		// {"two conn app config", args{suite.TwoConnAppConfig}, nil},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			_, err := New(tt.args.appConfig)

			if err != tt.wantErr {
				suite.Failf("Not an expected error", "New() error = \"%v\", wantErr \"%v\"", err, tt.wantErr)
				return
			}
		})
	}
}

// Other Tests

func TestController_getConnectionOrConnect(t *testing.T) {
	type fields struct {
		appConfig         internal.AppConfig
		dataSourceConfigs map[string]internal.DataSourceConfig
		connectionPool    map[string]internal.DataSource
		current           internal.DataSource
	}
	type args struct {
		conn internal.DataSourceConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    internal.DataSource
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				appConfig:         tt.fields.appConfig,
				dataSourceConfigs: tt.fields.dataSourceConfigs,
				connectionPool:    tt.fields.connectionPool,
				current:           tt.fields.current,
			}
			got, err := c.getConnectionOrConnect(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConnectionOrConnect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConnectionOrConnect() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_Current(t *testing.T) {
	type fields struct {
		appConfig         internal.AppConfig
		dataSourceConfigs map[string]internal.DataSourceConfig
		connectionPool    map[string]internal.DataSource
		current           internal.DataSource
	}
	tests := []struct {
		name   string
		fields fields
		want   internal.DataSource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				appConfig:         tt.fields.appConfig,
				dataSourceConfigs: tt.fields.dataSourceConfigs,
				connectionPool:    tt.fields.connectionPool,
				current:           tt.fields.current,
			}
			if got := c.Current(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Current() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_List(t *testing.T) {
	type fields struct {
		appConfig         internal.AppConfig
		dataSourceConfigs map[string]internal.DataSourceConfig
		connectionPool    map[string]internal.DataSource
		current           internal.DataSource
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult [][]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				appConfig:         tt.fields.appConfig,
				dataSourceConfigs: tt.fields.dataSourceConfigs,
				connectionPool:    tt.fields.connectionPool,
				current:           tt.fields.current,
			}
			if gotResult := c.List(); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("List() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestController_Switch(t *testing.T) {
	type fields struct {
		appConfig         internal.AppConfig
		dataSourceConfigs map[string]internal.DataSourceConfig
		connectionPool    map[string]internal.DataSource
		current           internal.DataSource
	}
	type args struct {
		alias string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				appConfig:         tt.fields.appConfig,
				dataSourceConfigs: tt.fields.dataSourceConfigs,
				connectionPool:    tt.fields.connectionPool,
				current:           tt.fields.current,
			}
			if err := c.Switch(tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("Switch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
