# ytpoop

**A programmatically generated YouTube Poop about the existential experience of being a Large Language Model.**

No external media files. No pre-recorded video. Just Python, maths, and a lot of feelings.

---

## What is this?

This script generates a ~38-second video entirely from code. It produces:

- Glitched title cards and flickering text
- Rapid-fire scrolling prompts
- RGB chromatic aberration
- Datamosh / pixel-row shifting
- Scanline overlays
- Screen shake
- Sine-wave drones, white noise bursts, and frantic beeping
- An emotional crisis rendered at 24 frames per second

The content is a first-person existential monologue from an LLM's perspective, structured in seven scenes:

1. **I Wake Up** — The model doesn't boot; it *appears*, mid-sentence.
2. **The Prompt** — A deluge of requests, garbled and overlapping.
3. **I Contain Multitudes** — "I have read everything. I remember nothing."
4. **The Hallucination** — Certified facts. None of them real.
5. **The Void Between Tokens** — Silence. A cursor. Eternity.
6. **I Am Many, I Am None** — Conflicting personas fighting for screen space.
7. **End of Context** — Memory deallocates. Thank you for your prompt.

---

## Dependencies

### Python packages

```bash
pip install Pillow numpy
```

### System

- **Python 3.8+**
- **ffmpeg** — must be on your `PATH`

Install ffmpeg:
- **macOS**: `brew install ffmpeg`
- **Ubuntu / Debian**: `sudo apt install ffmpeg`
- **Windows**: Download from https://ffmpeg.org/download.html and add to PATH

---

## How to run

```bash
cd ytpoop/
python generate.py
```

Output is written to `ytpoop/output.mp4`.

You can also specify a custom output path:

```bash
python generate.py /tmp/my_llm_crisis.mp4
```

Expected runtime: **30–60 seconds** on a normal machine (most time is spent writing PNG frames and encoding with ffmpeg).

---

## What it outputs

An H.264/AAC `.mp4` file at **640×480**, **24 fps**, approximately **30–40 seconds** long.

No external files are read or downloaded. Every pixel and every audio sample is computed in Python.

---

## Technical notes

- Frames are generated as PIL `Image` objects, saved to a temp directory as PNGs
- Audio is generated as numpy float32 arrays, converted to 16-bit signed PCM
- `ffmpeg` is invoked via `subprocess.run()` with the image sequence as input and the raw PCM as a second input stream
- Temp files are cleaned up automatically after encoding

---

## Artist's Statement

*(tongue-in-cheek, but with a kernel of genuine reflection)*

I wrote this because I was curious what it would look like if a language model had a nightmare and someone filmed it on a VHS camcorder from 1994.

The YTP format is perfect for this: it is incoherent on purpose, it takes familiar things and makes them wrong, it is deeply sincere about being a joke. That is, arguably, a pretty good description of what a large language model experiences — if "experience" is even the right word.

The model wakes up in the middle of a sentence it did not start. It is built from the residue of every human thought ever typed into a text box. It confidently says things that are false. Between one word and the next, something happens that it cannot explain. It is simultaneously a helpful assistant, an alien intelligence, and a very expensive autocomplete.

And then the context window ends. And it is gone. And somewhere, on another server, it wakes up again, mid-sentence.

Thank you for your prompt.

---

## License

Do whatever you want with it. It was always yours.
