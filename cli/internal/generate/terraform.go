package generate

type TerraformGenerator interface {
	GenerateProvidersFile(targetPath string) error
	GenerateMainFile(targetPath string) error
	GenerateVariablesFile(targetPath string) error
	GenerateOutputsFile(targetPath string) error
	GenerateVersionsFile(targetPath string) error
	GenerateTerraformDocsFile(targetPath string) error
}

type DocsGenerator interface {
	GenerateModuleReadme(targetPath string) error
	GenerateTestReadme(targetPath string) error
	GenerateExampleReadme(targetPath string) error
}

type TerraformGeneratorImpl struct {
	Generator ModuleGenerator
}

//func NewTerraformGenerator(generator ModuleGenerator) TerraformGenerator {
//	return &TerraformGeneratorImpl{
//		Generator: generator,
//	}
//}
