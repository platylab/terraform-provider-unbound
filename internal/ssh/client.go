package ssh

import (
  "encoding/json"
  "errors"
  "os"
  "strings"

  "golang.org/x/crypto/ssh"
)

// SSHClient wraps the SSH client
type SSHClient struct {
  client *ssh.Client
}

// NewSSHClient initializes a new SSH client
func NewSSHClient(host, username, private_key_path string) (*SSHClient, error) {
  // Load the private key
  key, err := os.ReadFile(private_key_path)
  if err != nil {
    return nil, errors.New("failed to read private key: " + err.Error())
  }

  // Create the signer for the private key
  signer, err := ssh.ParsePrivateKey(key)
  if err != nil {
    return nil, errors.New("failed to parse private key: " + err.Error())
  }

  // Set up the SSH client configuration
  config := &ssh.ClientConfig{
    User: username,
    Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
    HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use with caution
  }

  // Establish the connection
  client, err := ssh.Dial("tcp", host, config)
  if err != nil {
    return nil, errors.New("failed to dial SSH connection: " + err.Error())
  }
  return &SSHClient{client: client}, nil
}

// ExecuteCommand runs a command on the remote host
func (s *SSHClient) ExecuteCommand(command string) (string, error) {
  // Create a new session for each command
  session, err := s.client.NewSession()
  if err != nil {
    return "", errors.New("failed to create SSH session: " + err.Error())
  }
  defer session.Close() // Ensure the session is closed after execution

  output, err := session.CombinedOutput(command)
  if err != nil {
    return "", errors.New("SSH command failed: " + err.Error())
  }

  return string(output), nil
}

// ExecuteEndpoint runs a script on the remote host and parses its JSON output.
func (s *SSHClient) ExecuteEndpoint(command string, args map[string]string) (map[string]interface{}, error) {
  // Use a default executable if not provided
  var execPath string
  if strings.Contains(command, " ") { // If the command contains spaces, assume it's a full path with subcommand
    parts := strings.Split(command, " ")
    execPath = parts[0]
    command = strings.Join(parts[1:], " ")
  } else {
    execPath = "/usr/local/bin/unbound-endpoint"
  }

  // Build command arguments
  var cmdArgs []string
  cmdArgs = append(cmdArgs, execPath)
  cmdArgs = append(cmdArgs, command)
  for key, value := range args {
    cmdArgs = append(cmdArgs, "--"+key, value)
  }
  
  // Execute the command and capture output
  output, err := s.ExecuteCommand(strings.Join(cmdArgs, " "))
  if err != nil {
    return nil, errors.New("failed to execute script: " + err.Error())
  }

  // Parse JSON output
  var result map[string]interface{}
  err = json.Unmarshal([]byte(output), &result)
  if err != nil {
    return nil, errors.New("failed to parse JSON output: " + err.Error())
  }

  return result, nil
}

// Close cleans up the SSH client
func (s *SSHClient) Close() error {
  return s.client.Close()
}
