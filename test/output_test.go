package test

import (
	"lang/utils"
	"testing"
)

func TestCase001(t *testing.T) {
	utils.AssertProgramOutput("testcases/001.lang", "1\n", t)
}

func TestCase002(t *testing.T) {
	utils.AssertProgramOutput("testcases/002.lang", "7\n", t)
}

func TestCase003(t *testing.T) {
	utils.AssertProgramOutput("testcases/003.lang", "15\n", t)
}

func TestCase004(t *testing.T) {
	utils.AssertProgramOutput("testcases/004.lang", "12\n", t)
}

func TestCase005(t *testing.T) {
	utils.AssertProgramOutput("testcases/005.lang", "", t)
}

func TestCase006(t *testing.T) {
	utils.AssertProgramOutput("testcases/006.lang", "", t)
}

func TestCase007(t *testing.T) {
	utils.AssertCompilerFails("testcases/007.lang", t)
}

func TestCase008(t *testing.T) {
	utils.AssertCompilerFails("testcases/008.lang", t)
}

func TestCase009(t *testing.T) {
	utils.AssertCompilerFails("testcases/009.lang", t)
}

func TestCase010(t *testing.T) {
	utils.AssertCompilerFails("testcases/010.lang", t)
}

func TestCase011(t *testing.T) {
	utils.AssertProgramOutput("testcases/011.lang", "5\n", t)
}

func TestCase012(t *testing.T) {
	utils.AssertCompilerFails("testcases/012.lang", t)
}

func TestCase013(t *testing.T) {
	utils.AssertProgramOutput("testcases/013.lang", "3\n2\n", t)
}

func TestCase014(t *testing.T) {
	utils.AssertProgramOutput("testcases/014.lang", "3\n2\n1\n0\n", t)
}

func TestCase015(t *testing.T) {
	utils.AssertProgramOutput("testcases/015.lang", "42\n", t)
}

func TestCase016(t *testing.T) {
	utils.AssertProgramOutput("testcases/016.lang", "42\n42\n", t)
}

func TestCase017(t *testing.T) {
	utils.AssertProgramOutput("testcases/017.lang", "50\n", t)
}

func TestCase018(t *testing.T) {
	utils.AssertProgramOutput("testcases/018.lang", "48\n", t)
}

func TestCase019(t *testing.T) {
	utils.AssertProgramOutput("testcases/019.lang", "5\n", t)
}

func TestCase020(t *testing.T) {
	utils.AssertProgramOutput("testcases/020.lang", "21\n", t)
}

func TestCase021(t *testing.T) {
	utils.AssertProgramOutput("testcases/021.lang", "2\n", t)
}

func TestCase022(t *testing.T) {
	utils.AssertProgramOutput("testcases/022.lang", "12\n", t)
}

func TestCase023(t *testing.T) {
	utils.AssertProgramOutput("testcases/023.lang", "230\n", t)
}

func TestCase024(t *testing.T) {
	utils.AssertCompilerFails("testcases/024.lang", t)
}

func TestCase025(t *testing.T) {
	utils.AssertCompilerFails("testcases/025.lang", t)
}

func TestCase026(t *testing.T) {
	utils.AssertProgramOutput("testcases/026.lang", "10\n", t)
}

func TestCase027(t *testing.T) {
	utils.AssertProgramOutput("testcases/027.lang", "6\n", t)
}

func TestCase028(t *testing.T) {
	utils.AssertProgramOutput("testcases/028.lang", "32\n", t)
}

func TestCase029(t *testing.T) {
	utils.AssertProgramOutput("testcases/029.lang", "16\n", t)
}

func TestCase030(t *testing.T) {
	utils.AssertProgramOutput("testcases/030.lang", "7\n", t)
}

func TestCase031(t *testing.T) {
	utils.AssertProgramOutput("testcases/031.lang", "20\n", t)
}

func TestCase032(t *testing.T) {
	utils.AssertProgramOutput("testcases/032.lang", "1\n", t)
}

func TestCase033(t *testing.T) {
	utils.AssertProgramOutput("testcases/033.lang", "0\n", t)
}

func TestCase034(t *testing.T) {
	utils.AssertProgramOutput("testcases/034.lang", "1\n", t)
}

func TestCase035(t *testing.T) {
	utils.AssertProgramOutput("testcases/035.lang", "1\n", t)
}

func TestCase036(t *testing.T) {
	utils.AssertCompilerFails("testcases/036.lang", t)
}

func TestCase037(t *testing.T) {
	utils.AssertProgramOutput("testcases/037.lang", "1\n0\n0\n0\n", t)
}

func TestCase038(t *testing.T) {
	utils.AssertProgramOutput("testcases/038.lang", "3\n", t)
}

func TestCase039(t *testing.T) {
	utils.AssertProgramOutput("testcases/039.lang", "0\n1\n", t)
}

func TestCase040(t *testing.T) {
	utils.AssertProgramOutput("testcases/040.lang", "-5\n-4\n-3\n", t)
}

func TestCase041(t *testing.T) {
	utils.AssertProgramOutput("testcases/041.lang", "3\n-3\n", t)
}

func TestCase042(t *testing.T) {
	utils.AssertProgramOutput("testcases/042.lang", "24\n", t)
}

func TestCase043(t *testing.T) {
	utils.AssertProgramOutput("testcases/043.lang", "0\n2\n1\n", t)
}

func TestCase044(t *testing.T) {
	utils.AssertProgramCrashes("testcases/044.lang", t)
}

func TestCase045(t *testing.T) {
	utils.AssertProgramOutput("testcases/045.lang", "-1\n-2\n", t)
}
