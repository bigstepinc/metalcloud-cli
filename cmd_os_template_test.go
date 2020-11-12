package main

import (
	"testing"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	mock_metalcloud "github.com/bigstepinc/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestOSTemplatesListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	list := map[string]metalcloud.OSTemplate{
		"test": {
			VolumeTemplateID:        10,
			VolumeTemplateLabel:     "test",
			OSAssetBootloaderOSBoot: 100,
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		OSTemplates().
		Return(&list, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	client.EXPECT().
		OSAssetGet(list["test"].OSAssetBootloaderOSBoot).
		Return(&asset, nil).
		AnyTimes()

	//test json

	expectedFirstRow := map[string]interface{}{
		"ID":    10,
		"LABEL": "test",
	}

	testListCommand(templatesListCmd, nil, client, expectedFirstRow, t)

}

func TestOSTemplateCreateCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:        10,
		VolumeTemplateLabel:     "test",
		OSAssetBootloaderOSBoot: 100,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		OSTemplateCreate(gomock.Any()).
		Return(&tmpl, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	client.EXPECT().
		OSAssetGet(tmpl.OSAssetBootloaderOSBoot).
		Return(&asset, nil).
		AnyTimes()

	//test json

	cases := []CommandTestCase{
		{
			name: "create-good1",
			cmd: MakeCommand(map[string]interface{}{
				"label":            tmpl.VolumeTemplateLabel,
				"display_name":     "test",
				"boot_type":        "pxe",
				"os_type":          "centos",
				"os_version":       "t0",
				"os_architecture":  "t0",
				"initial_user":     "t0",
				"initial_password": "t0",
				"initial_ssh_port": 22,
			}),
			good: true,
		},
		{
			name: "with use-autogenerated-initial-password ",
			cmd: MakeCommand(map[string]interface{}{
				"label":                              tmpl.VolumeTemplateLabel,
				"display_name":                       "test",
				"boot_type":                          "pxe",
				"os_type":                            "centos",
				"os_version":                         "t0",
				"os_architecture":                    "t0",
				"initial_user":                       "t0",
				"use_autogenerated_initial_password": true,
				"initial_ssh_port":                   22,
			}),
			good: true,
		},
		{
			name: "missing label",
			cmd: MakeCommand(map[string]interface{}{
				//"label":            tmpl.VolumeTemplateLabel,
				"display_name":     "test",
				"boot_type":        "pxe",
				"os_type":          "centos",
				"os_version":       "t0",
				"os_architecture":  "t0",
				"initial_user":     "t0",
				"initial_password": "t0",
				"initial_ssh_port": 22,
			}),
			good: false,
		},
		{
			name: "missing initial_password",
			cmd: MakeCommand(map[string]interface{}{
				"label":           tmpl.VolumeTemplateLabel,
				"display_name":    "test",
				"boot_type":       "pxe",
				"os_type":         "centos",
				"os_version":      "t0",
				"os_architecture": "t0",
				"initial_user":    "t0",
				//	"initial_password": "t0",
				"initial_ssh_port": 22,
			}),
			good: false,
		},
		{
			name: "missing initial_user",
			cmd: MakeCommand(map[string]interface{}{
				"label":           tmpl.VolumeTemplateLabel,
				"display_name":    "test",
				"boot_type":       "pxe",
				"os_type":         "centos",
				"os_version":      "t0",
				"os_architecture": "t0",
				//"initial_user":    "t0",
				"initial_password": "t0",
				"initial_ssh_port": 22,
			}),
			good: false,
		},
		{
			name: "missing either",
			cmd: MakeCommand(map[string]interface{}{
				"label":           tmpl.VolumeTemplateLabel,
				"display_name":    "test",
				"boot_type":       "pxe",
				"os_type":         "centos",
				"os_version":      "t0",
				"os_architecture": "t0",
				//"initial_user":    "t0",
				//"initial_password": "t0",
				"initial_ssh_port": 22,
			}),
			good: false,
		},
		{
			name: "both password options",
			cmd: MakeCommand(map[string]interface{}{
				"label":                              tmpl.VolumeTemplateLabel,
				"display_name":                       "test",
				"boot_type":                          "pxe",
				"os_type":                            "centos",
				"os_version":                         "t0",
				"os_architecture":                    "t0",
				"initial_user":                       "t0",
				"initial_password":                   "t0",
				"use_autogenerated_initial_password": true,
				"initial_ssh_port":                   22,
			}),
			good: false,
		},
	}

	testCreateCommand(templateCreateCmd, cases, client, t)

}

func TestOSTemplateMakePrivateCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"t1": tmpl,
	}

	user := metalcloud.User{
		UserID: 1,
	}

	user1 := metalcloud.User{
		UserEmail: "test",
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		AnyTimes()

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		MinTimes(1)

	client.EXPECT().
		UserGet(gomock.Any()).
		Return(&user, nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(gomock.Any()).
		Return(&user1, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateMakePrivate(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": 10,
				"user_id":             1,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good2",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": "test",
				"user_id":             1,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good3",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": 10,
				"user_id":             "test",
			}),
			good: true,
			id:   0,
		},
		{
			name: "template not found",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": "test1",
				"user_id":             1,
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing template id or name",
			cmd: MakeCommand(map[string]interface{}{
				"user_id": 1,
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing user id or email",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": "test",
			}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(templateMakePrivateCmd, cases, client, t)
}

func TestOSTemplateMakePublicCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"t1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		AnyTimes()

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateMakePublic(gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": 10,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good2",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": "test",
			}),
			good: true,
			id:   0,
		},
		{
			name: "template not found",
			cmd: MakeCommand(map[string]interface{}{
				"template_id_or_name": "test1",
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing template id or name",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(templateMakePublicCmd, cases, client, t)
}
