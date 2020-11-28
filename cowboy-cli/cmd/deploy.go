/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var region, image, name, ports, replica string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		address := os.Getenv("API_URL")
		if os.Getenv("API_URL") == "" {
			log.Fatal("Please add api url")
		}
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to connect to %s: %v", address, err)
		}
		c := pb.NewPlatformServiceClient(conn)

		// Create metadata and context.
		fmt.Println(os.Getenv("COWBOY_TOKEN"))
		if os.Getenv("COWBOY_TOKEN") == "" {
			log.Fatal("Please add token")
		}
		md := metadata.Pairs("authorization", "Bearer "+os.Getenv("COWBOY_TOKEN"))
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		var header, trailer metadata.MD
		r, err := c.CreateApplication(ctx, &pb.CreateApplicationRequest{
			Image:   image,
			Ports:   ports,
			Replica: replica,
			Region:  region,
			Name:    name,
		}, grpc.Header(&header), grpc.Trailer(&trailer))

		if err != nil {
			log.Fatal(err)
		}
		log.Println("Error:", r.Error)
		log.Println("Message:", r.Message)
		log.Println("Data:", r.Data)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.PersistentFlags().StringVarP(&image, "image", "i", "", "container Image")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region (required)")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "name of test")
	rootCmd.PersistentFlags().StringVarP(&ports, "ports", "p", "", "port of container")
	rootCmd.PersistentFlags().StringVarP(&replica, "replica", "c", "2", "replica count")

}
