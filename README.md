# Embedded Reflash Tool

A small utility built primarily to replace AMI MegaRAC with OpenBMC in the field. This tool can flash firmware images directly to MTD devices and includes safety measures like remounting filesystems read-only before flashing.

## How It Works

The tool performs the following steps:

1. **Safety Check**: Ensures it's running as root
2. **Filesystem Protection**: Remounts read-write filesystems as read-only to prevent corruption during flash operations (can be skipped with `-skip-remount`)
3. **Image Preparation**: Either uses a built-in embedded image or reads from a specified file
4. **Flash Operation**: Uses the flashcp library to write the image to the specified MTD device
5. **System Reset**: Automatically reboots the system after successful flashing

## Usage

First, you need to get the compiled binary onto the target system that should be reflashed. Common methods include:

### Transfer Methods

**Via SCP:**
```bash
scp embedded_reflash root@target-ip:/tmp/
```

**Via SSH and wget:**
```bash
ssh root@target-ip
wget http://your-server/embedded_reflash -O /tmp/embedded_reflash
chmod +x /tmp/embedded_reflash
```

**Via existing web interface file upload (if available)**

### Running the Tool

Once the binary is on the target system:

**Basic usage with external image:**
```bash
./embedded_reflash -image /path/to/firmware.img -device /dev/mtd0
```

**Using built-in embedded image:**
```bash
./embedded_reflash -builtin -image /tmp/firmware.img -device /dev/mtd0
```

**Skip filesystem remounting (use with caution):**
```bash
./embedded_reflash -skip-remount -image /path/to/firmware.img
```

### Command Line Options

- `-device`: Flash device to write to (default: `/dev/mtd0`)
- `-image`: Flash image file path, or destination path when using `-builtin`
- `-builtin`: Use the image embedded in the binary
- `-skip-remount`: Skip remounting filesystems read-only (use with caution)

## Building

### Requirements

- Go programming language compiler (1.25 or later recommended)

### Build Steps

First, clone this repository and navigate to it:

```bash
git clone <repository-url>
cd embedded_reflash
```

**For builtin mode:** If you want to use the `-builtin` flag to embed an image directly in the binary, you must replace the empty `builtin.img` file with your actual firmware image before compiling:

```bash
cp /path/to/your/firmware.img builtin.img
```

Then build for Aspeed type SoCs:

**Pre Go 1.22:**

```bash
GOARCH=arm GOARM=5 go build -ldflags="-s -w"
```

**Go 1.22 or later:**

```bash
GOARCH=arm GOARM=6,softfloat go build -ldflags="-s -w" # AST2500
GOARCH=arm GOARM=7,softfloat go build -ldflags="-s -w" # AST2600
```

## Safety Notes

- Always ensure you have a recovery method available before flashing
- The tool will automatically reboot the system after flashing
- Filesystem remounting helps prevent corruption but may cause temporary service interruptions
- Only use `-skip-remount` if you understand the risks
