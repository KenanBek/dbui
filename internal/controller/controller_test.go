package controller

import (
	"dbui/internal"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

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

func TestNew(t *testing.T) {
	type args struct {
		appConfig internal.AppConfig
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	emptyAppConfig := NewMockAppConfig(ctrl)
	emptyAppConfig.EXPECT().DataSourceConfigs().Return(map[string]internal.DataSourceConfig{})

	tests := []struct {
		name    string
		args    args
		wantC   *Controller
		wantErr bool
	}{
		{"empty appConfig", args{emptyAppConfig}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := New(tt.args.appConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("New() gotC = %v, want %v", gotC, tt.wantC)
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
