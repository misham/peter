package extractor_test

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/misham/peter/extractor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Extract", func() {

	var extractionDest string

	BeforeEach(func() {
		var err error

		extractionDest, err = ioutil.TempDir("", "extractor")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(extractionDest)
	})

	Describe("with bad file that is", func() {
		Context("missing", func() {
			It("returns an error", func() {
				file := "unknown_file"
				err := extractor.Extract(&extractionDest, &file)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("unknown format", func() {
			It("returns an error", func() {
				file := "fixtures/uncompressed_file.txt"
				err := extractor.Extract(&extractionDest, &file)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("with good file that is", func() {
		Context("a zip archive", func() {
			Context("multiple files", func() {
				It("extracts them successfully", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/zip/multiple_files.zip")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					Expect(path.Join(extractionDest, "test1.txt")).To(BeARegularFile())
					Expect(path.Join(extractionDest, "test2.txt")).To(BeARegularFile())
				})
			})

			Context("nested directories", func() {
				It("extracts them successfully", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/zip/nested_dir.zip")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					Expect(path.Join(extractionDest, "nested")).To(BeADirectory())
					Expect(path.Join(extractionDest, "nested/test.txt")).To(BeARegularFile())
				})
			})

			Context("permissions", func() {
				It("are preserved", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/zip/perms.zip")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					testFile := path.Join(extractionDest, "test.sh")
					Expect(testFile).To(BeARegularFile())

					fileInfo, err := os.Stat(testFile)
					Expect(err).NotTo(HaveOccurred())

					Expect(fileInfo.Mode().Perm()).To(BeEquivalentTo(0755))
				})
			})

			Context("symlinks", func() {
				It("are created correctly", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/zip/symlink.zip")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					Expect(path.Join(extractionDest, "file.txt")).To(BeARegularFile())
					Expect(path.Join(extractionDest, "link.txt")).To(BeAnExistingFile())

					symlinkInfo, err := os.Lstat(path.Join(extractionDest, "link.txt"))
					Expect(err).NotTo(HaveOccurred())
					Expect(symlinkInfo.Mode() & 0755).To(Equal(os.FileMode(0755)))

					target, err := os.Readlink(path.Join(extractionDest, "link.txt"))
					Expect(err).NotTo(HaveOccurred())
					Expect(target).To(Equal("file.txt"))
				})
			})
		})

		Context("a tgz archive", func() {
			Context("multiple files", func() {
				It("extracts them successfully", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/tgz/multiple_files.tgz")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					Expect(path.Join(extractionDest, "test1.txt")).To(BeARegularFile())
					Expect(path.Join(extractionDest, "test2.txt")).To(BeARegularFile())
				})
			})

			Context("nested directories", func() {
				It("extracts them successfully", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/tgz/nested_dir.tgz")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					Expect(path.Join(extractionDest, "nested")).To(BeADirectory())
					Expect(path.Join(extractionDest, "nested/test.txt")).To(BeARegularFile())
				})
			})

			Context("permissions", func() {
				It("are preserved", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/tgz/perms.tgz")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					testFile := path.Join(extractionDest, "test.sh")
					Expect(testFile).To(BeARegularFile())

					fileInfo, err := os.Stat(testFile)
					Expect(err).NotTo(HaveOccurred())

					Expect(fileInfo.Mode().Perm()).To(BeEquivalentTo(0755))
				})
			})

			Context("symlinks", func() {
				It("are created correctly", func() {
					cwd, err := os.Getwd()
					Expect(err).NotTo(HaveOccurred())

					file := path.Join(cwd, "fixtures/tgz/symlink.tgz")
					err = extractor.Extract(&extractionDest, &file)
					Expect(err).NotTo(HaveOccurred())

					Expect(path.Join(extractionDest, "file.txt")).To(BeARegularFile())
					Expect(path.Join(extractionDest, "link.txt")).To(BeAnExistingFile())

					symlinkInfo, err := os.Lstat(path.Join(extractionDest, "link.txt"))
					Expect(err).NotTo(HaveOccurred())
					Expect(symlinkInfo.Mode() & 0755).To(Equal(os.FileMode(0755)))

					target, err := os.Readlink(path.Join(extractionDest, "link.txt"))
					Expect(err).NotTo(HaveOccurred())
					Expect(target).To(Equal("file.txt"))
				})
			})
		})
	})
})
