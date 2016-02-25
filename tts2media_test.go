package tts2media

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var (
	picowav       *PicoTTSSpeech
	videoFilename string
)

// TestPreload calls SetDataDir to set the output directories
func TestPreload(t *testing.T) {
	webmB64 := "GkXfo0AgQoaBAUL3gQFC8oEEQvOBCEKCQAR3ZWJtQoeBAkKFgQIYU4BnQI0VSalmQCgq17FAAw9CQE2AQAZ3aGFtbXlXQUAGd2hhbW15RIlACECPQAAAAAAAFlSua0AxrkAu14EBY8WBAZyBACK1nEADdW5khkAFVl9WUDglhohAA1ZQOIOBAeBABrCBCLqBCB9DtnVAIueBAKNAHIEAAIAwAQCdASoIAAgAAUAmJaQAA3AA/vz0AAA="
	mp4B64 := "AAAAHGZ0eXBpc29tAAACAGlzb21pc28ybXA0MQAAAAhmcmVlAAAAGm1kYXQAAAGzABAHAAABthBgUYI9t+8AAAMNbW9vdgAAAGxtdmhkAAAAAMXMvvrFzL76AAAD6AAAACoAAQAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAABhpb2RzAAAAABCAgIAHAE/////+/wAAAiF0cmFrAAAAXHRraGQAAAAPxcy++sXMvvoAAAABAAAAAAAAACoAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAABAAAAAAAgAAAAIAAAAAAG9bWRpYQAAACBtZGhkAAAAAMXMvvrFzL76AAAAGAAAAAEVxwAAAAAALWhkbHIAAAAAAAAAAHZpZGUAAAAAAAAAAAAAAABWaWRlb0hhbmRsZXIAAAABaG1pbmYAAAAUdm1oZAAAAAEAAAAAAAAAAAAAACRkaW5mAAAAHGRyZWYAAAAAAAAAAQAAAAx1cmwgAAAAAQAAAShzdGJsAAAAxHN0c2QAAAAAAAAAAQAAALRtcDR2AAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAgACABIAAAASAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGP//AAAAXmVzZHMAAAAAA4CAgE0AAQAEgICAPyARAAAAAAMNQAAAAAAFgICALQAAAbABAAABtYkTAAABAAAAASAAxI2IAMUARAEUQwAAAbJMYXZjNTMuMzUuMAaAgIABAgAAABhzdHRzAAAAAAAAAAEAAAABAAAAAQAAABxzdHNjAAAAAAAAAAEAAAABAAAAAQAAAAEAAAAUc3RzegAAAAAAAAASAAAAAQAAABRzdGNvAAAAAAAAAAEAAAAsAAAAYHVkdGEAAABYbWV0YQAAAAAAAAAhaGRscgAAAAAAAAAAbWRpcmFwcGwAAAAAAAAAAAAAAAAraWxzdAAAACOpdG9vAAAAG2RhdGEAAAABAAAAAExhdmY1My4yMS4x"
	pngB64 := "iVBORw0KGgoAAAANSUhEUgAAAA0AAAANCAMAAABFNRROAAAAB3RJTUUH1wQVFxccWUPqdgAAAAlwSFlzAAALEQAACxEBf2RfkQAAAARnQU1BAACxjwv8YQUAAAAVUExURQAAAIbB4BVEcGuk2kqFxjhsn0J6tlg7/OkAAAABdFJOUwBA5thmAAAAQ0lEQVR42lWNQRIAMAQDlcj/n1zUId0LawKz4ZjCTx0ixymWvtFoPImuPS9IVu1lcHlbMX3sDaRaugStLnYc+3pCKLtcswFHkVxAmgAAAABJRU5ErkJggg=="

	saveDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal("Error creating a temporary directory: " + err.Error())
	}

	saveDir += "/"
	if err = os.Mkdir(saveDir+"tmp", 0700); err != nil {
		t.Fatal("Error creating " + saveDir + "tmp")
	}

	SetDataDir(saveDir)

	// writes file to tempPath name, as payload.ext
	writeMedia := func(m string, ext string) {
		stream, err := base64.StdEncoding.DecodeString(m)
		if err != nil {
			t.Fatal("Error decoding base64 encoded file, payload." + ext)
		}
		if err := ioutil.WriteFile(tempPath+"payload."+ext, stream, 0700); err != nil {
			t.Fatal("Error creating " + tempPath + "payload." + ext)
		}
	}

	writeMedia(pngB64, "png")
	writeMedia(mp4B64, "mp4")
	writeMedia(webmB64, "webm")
}

// TestEspeak checks that NewEspeakSpeech() creates the WAV files,
// and that the file can be removed by media.RemoveWAV()
func TestEspeak(t *testing.T) {
	espwav := &EspeakSpeech{"This is some sample text for testing text to speech engines.",
		"en", "135", "m", "0", "high", "50"}

	media, err := espwav.NewEspeakSpeech()
	if err != nil {
		t.Fatal("NewEspeakSpeech returned error:", err)
	}

	filename := dataPath + media.Filename + ".wav"

	// t.Log("type =", reflect.TypeOf(espwav))
	t.Log("filename =", filename)
	if _, err = os.Stat(filename); err != nil {
		t.Fatal("There was an error opening the WAV file:", err)
	}

	media.RemoveWAV()
	filename = dataPath + media.Filename + ".wav"
	if _, err = os.Stat(filename); err == nil {
		t.Fatal("There was an error removing the WAV file:", err)
	}

	return
}

// TestEspeak checks that media.ToAudio() creates the MP3 and OGG files after creating
// the WAV files, and that the WAV file can be removed by media.RemoveWAV()
// This function also tests media.ImageToVideo() and sets the filename for other tests
func TestEspeakToAudio(t *testing.T) {
	espwav := &EspeakSpeech{"This is some sample text for testing text to speech engines.",
		"en", "135", "m", "0", "medium", "50"}
	media, err := espwav.NewEspeakSpeech()
	if err != nil {
		t.Fatal("NewEspeakSpeech returned error:", err)
	}

	wavFile := dataPath + media.Filename + ".wav"
	t.Log(wavFile)

	err = media.ToAudio()
	if err != nil {
		t.Fatal("There was an error creating the audio files:", err)
	}

	filename := dataPath + media.Filename

	t.Log("filename =", filename+".mp3")
	if _, err = os.Stat(filename + ".mp3"); err != nil {
		t.Fatal("There was an error opening the audio files:", err)
	}

	t.Log("filename =", filename+".ogg")
	if _, err = os.Stat(filename + ".ogg"); err != nil {
		t.Fatal("There was an error opening the audio files:", err)
	}

	tempFile := tempPath + media.Filename

	ext := "png"
	if !copy(tempPath+"payload."+ext, tempFile+"."+ext) {
		t.Fatal("Could not copy "+tempPath+"payload.png to", tempFile)
	}

	if videoFilename, err = ImageToVideo(media.Filename, ext); err != nil {
		t.Fatal("There was an error in ImageToVideo()", err)
	}

	os.Remove(tempFile + "." + ext)
	os.Remove(tempFile)

	media.RemoveWAV()
	t.Log("filename =", filename+".wav")
	if _, err = os.Stat(filename + ".wav"); err == nil {
		t.Fatal("There was an error removing the WAV file:", err)
	}

	os.Remove(filename + ".mp3")
	os.Remove(filename + ".ogg")
	os.Remove(dataPath + videoFilename + ".webm")
	os.Remove(dataPath + videoFilename + ".mp4")
}

// TestPicoTTS checks that NewPicoTTSSpeech() creates the WAV files,
// and that the file can be removed by media.RemoveWAV()
func TestPicoTTS(t *testing.T) {
	picowav := &PicoTTSSpeech{"This is some sample text for testing text to speech engines.",
		"en-US", "medium"}
	media, err := picowav.NewPicoTTSSpeech()
	if err != nil {
		t.Fatal("NewPicoTTSSpeech returned error:", err)
	}

	filename := dataPath + media.Filename + ".wav"

	t.Log("filename =", filename)
	if _, err = os.Stat(filename); err != nil {
		t.Fatal("There was an error opening the WAV file:", err)
	}

	media.RemoveWAV()
	filename = dataPath + media.Filename + ".wav"
	if _, err = os.Stat(filename); err == nil {
		t.Fatal("There was an error removing the WAV file:", err)
	}

	return
}

// TestPicoTTSToAudio checks that media.ToAudio() creates the MP3 and OGG files after creating the
// WAV files, and that the WAV file can be removed by media.RemoveWAV()
func TestPicoTTSToAudio(t *testing.T) {
	picowav := &PicoTTSSpeech{"This is some sample text for testing text to speech engines." +
		"This is some sample text for testing text to speech engines. This is some sample" +
		" text for testing text to speech engines. This is some sample text for testing " +
		"text to speech engines. This is some sample text for testing text to speech " +
		"engines. This is some sample text for testing text to speech engines. This is " +
		"some sample text for testing text to speech engines. This is some sample text " +
		"for testing text to speech engines. This is some sample text for testing text to" +
		" speech engines. This is some sample text for testing text to speech engines. " +
		"This is some sample text for testing text to speech engines.", "en-US", "low"}

	media, err := picowav.NewPicoTTSSpeech()
	if err != nil {
		t.Fatal("NewPicoTTSSpeech returned error:", err)
	}

	wavFile := dataPath + media.Filename + ".wav"
	t.Log(wavFile)

	err = media.ToAudio()
	if err != nil {
		t.Fatal("There was an error creating the audio file:", err)
	}

	filename := dataPath + media.Filename

	t.Log("filename =", filename+".mp3")
	if _, err = os.Stat(filename + ".mp3"); err != nil {
		t.Fatal("There was an error opening the audio files:", err)
	}

	t.Log("filename =", filename+".ogg")
	if _, err = os.Stat(filename + ".ogg"); err != nil {
		t.Fatal("There was an error opening the audio files:", err)
	}

	media.RemoveWAV()
	t.Log("filename =", filename+".wav")
	if _, err = os.Stat(filename + ".wav"); err == nil {
		t.Fatal("There was an error removing the WAV file:", err)
	}

	if !copy(tempPath+"payload.mp4", tempPath+media.Filename) {
		t.Fatal("Could not copy " + tempPath + "payload.mp4 to " + tempPath + media.Filename)
	}

	duration, err := Duration(tempPath + media.Filename)
	if err != nil {
		t.Fatal("There was an error in Duration():", err)
	}
	t.Log("duration =", duration)

	if videoFilename, err = FromVideo(media.Filename, false); err != nil {
		t.Fatal("There was an error in FromVideo():", err)
	}

	os.Remove(filename + ".mp3")
	os.Remove(filename + ".ogg")
	os.Remove(filename + ".webm")
	os.Remove(filename + ".mp4")

	os.Remove(tempPath + media.Filename)

	os.Remove(dataPath + videoFilename + ".mp4")
	os.Remove(dataPath + videoFilename + ".webm")

	os.Remove(tempPath + "payload.webm")
	os.Remove(tempPath + "payload.png")
	os.Remove(tempPath + "payload.mp4")

	os.Remove(tempPath)
	os.Remove(dataPath)
}

func copy(in, out string) bool {

	r, err := os.Open(in)
	if err != nil {
		return false
	}
	defer r.Close()

	w, err := os.Create(out)
	if err != nil {
		return false
	}
	defer w.Close()

	// do the actual work
	if _, err = io.Copy(w, r); err != nil {
		return false
	}
	return true
}
