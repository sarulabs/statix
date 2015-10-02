package helpers

import "testing"

func TestRewritePath(t *testing.T) {
	if RewritePath("test", "/absolute") != "/absolute" {
		t.Error("absolute paths should not be modified")
	}

	if RewritePath("test", "////absolute/") != "/absolute" {
		t.Error("absolute paths should be cleaned")
	}

	if RewritePath("test", "absolute") != "test/absolute" {
		t.Error("basePath should be added to relative paths")
	}

	if RewritePath("/test", "absolute") != "/test/absolute" {
		t.Error("basePath should be added to relative paths")
	}

	if RewritePath("./test////", "absolute////") != "test/absolute" {
		t.Error("relative paths should be cleaned")
	}
}

func TestAddFileSuffix(t *testing.T) {
	if AddFileSuffix("myfile", "SUFFIX") != "myfileSUFFIX" {
		t.Error("suffix should be added to filename")
	}

	if AddFileSuffix(".////myfile", "SUFFIX") != "myfileSUFFIX" {
		t.Error("path should be cleaned")
	}

	if AddFileSuffix("test/myfile.test", "SUFFIX") != "test/myfileSUFFIX.test" {
		t.Error("suffix should be added before file extension")
	}
}

func TestMD5(t *testing.T) {
	if MD5([]byte("testMd5")) != "ef8efa55f449e3727c4df433ce7744c5" {
		t.Error("md5 should return ef8efa55f449e3727c4df433ce7744c5")
	}
}
