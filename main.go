package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"
)

const (
	tmpDir          = "/tmp"
	namespace       = "http://schemas.android.com/apk/res/android"
	versionCodeAttr = "versionCode"
	versionNameAttr = "versionName"
)

type Config struct {
	versionCode int32
	versionName string
}

func main() {
	versionCode := flag.Uint("versionCode", 0, "The versionCode to set")
	versionName := flag.String("versionName", "", "The versionName to set")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(flag.CommandLine.Output(), "Error: File path is required.")
		flag.Usage()
		os.Exit(2)
	}
	config := &Config{
		versionCode: int32(*versionCode),
		versionName: *versionName,
	}

	path := flag.Arg(0)

	if strings.HasSuffix(path, ".apk") {
		updateApk(path, config)
	} else if strings.HasSuffix(path, ".aar") {
		updateAar(path, config)
	} else {
		updateManifest(path, config)
	}
}

func updateApk(path string, config *Config) {
	file, err := ioutil.TempFile(tmpDir, "*.aar")
	if err != nil {
		log.Fatalln("Failed creating temp file:", err)
	}
	defer os.Remove(file.Name())

	out, err := exec.Command("aapt2", "convert", "-o", file.Name(), "--output-format", "proto", path).CombinedOutput()
	if err != nil {
		log.Fatalln("Failed executing aapt2:", err, string(out))
	}

	updateAar(file.Name(), config)

	out, err = exec.Command("aapt2", "convert", "-o", path, "--output-format", "binary", file.Name()).CombinedOutput()
	if err != nil {
		log.Fatalln("Failed executing aapt2:", err, string(out))
	}
}

func updateAar(path string, config *Config) {
	manifest, err := ioutil.TempFile(tmpDir, "AndroidManifest.*.xml")
	if err != nil {
		log.Fatalln("Failed creating temp file:", err)
	}
	defer os.Remove(manifest.Name())

	extractFromZip(path, "AndroidManifest.xml", manifest)
	updateManifest(manifest.Name(), config)
	addToZip(path, "AndroidManifest.xml", manifest)
}

func addToZip(path string, name string, source *os.File) {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		log.Fatalln("Failed reading zip:", err)
	}
	defer zipReader.Close()

	tmpZip, err := ioutil.TempFile(tmpDir, "*.aar")
	if err != nil {
		log.Fatalln("Failed creating temp file:", err)
	}
	defer os.Remove(tmpZip.Name())

	targetZipWriter := zip.NewWriter(tmpZip)
	defer targetZipWriter.Close()

	for _, zipItem := range zipReader.File {
		header, err := zip.FileInfoHeader(zipItem.FileInfo())
		if err != nil {
			log.Fatalln("Failed creating header:", err)
		}
		header.Name = zipItem.Name
		targetItem, err := targetZipWriter.CreateHeader(header)
		if err != nil {
			log.Fatalln("Failed creating header:", err)
		}
		if zipItem.Name == name {
			source.Seek(0, 0)
			_, err = io.Copy(targetItem, source)
			if err != nil {
				log.Fatalln("Failed copying to zip:", err)
			}
		} else {
			zipItemReader, err := zipItem.Open()
			if err != nil {
				log.Fatalln("Failed reading zip:", err)
			}
			_, err = io.Copy(targetItem, zipItemReader)
			if err != nil {
				log.Fatalln("Failed copying to zip:", err)
			}
			zipItemReader.Close()
		}
	}

	f, err := os.Create(path)
	if err != nil {
		log.Fatalln("Failed opening file:", err)
	}
	defer f.Close()
	tmpZip.Seek(0, 0)
	io.Copy(f, tmpZip)

	f2, err := os.Create("fafa.zip")
	if err != nil {
		log.Fatalln("Failed opening file:", err)
	}
	defer f2.Close()
	tmpZip.Seek(0, 0)
	io.Copy(f2, tmpZip)
}

func extractFromZip(path string, name string, target *os.File) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	f := findFile(r, name)
	if f == nil {
		log.Fatalln(errors.New("file is missing"))
	}

	innerFile, err := f.Open()
	if err != nil {
		log.Fatalln("Failed opening zip file's AndroidManifest.xml:", err)
	}
	defer innerFile.Close()
	_, err = io.Copy(target, innerFile)
	if err != nil {
		log.Fatal(err)
	}
}

func findFile(r *zip.ReadCloser, name string) *zip.File {
	for _, f := range r.File {
		if f.Name != name {
			continue
		}
		return f
	}
	return nil
}

func updateManifest(path string, config *Config) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	xmlNode := &XmlNode{}
	if err := proto.Unmarshal(in, xmlNode); err != nil {
		log.Fatalln("Failed to parse manifest:", err)
	}
	for _, attr := range xmlNode.GetElement().GetAttribute() {
		if attr.GetNamespaceUri() != namespace {
			continue
		}
		switch attr.GetName() {
		case versionCodeAttr:
			if config.versionCode > 0 {
				println(attr.String())
				prim := attr.GetCompiledItem().GetPrim()
				if x, ok := prim.GetOneofValue().(*Primitive_IntDecimalValue); ok {
					x.IntDecimalValue = int32(config.versionCode)
				}
			}
		case versionNameAttr:
			if config.versionName != "" {
				attr.Value = config.versionName
			}
		}
	}

	out, err := proto.Marshal(xmlNode)
	if err != nil {
		log.Fatalln("Error marshalling XML:", err)
	}
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Fatalln("Error writing file:", err)
	}
}
