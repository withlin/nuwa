package controllers

import (
	"context"
	"fmt"
	nuwav1 "github.com/yametech/nuwa/api/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	nuwaRoomFlag    = "nuwa.io/room"
	nuwaCabinetFlag = "nuwa.io/cabinet"
	nuwaHostFlag    = "nuwa.io/host"
)

type coordinator struct {
	Client       client.Client
	Index        int
	Name         string               `json:"name,omitempty"`
	Coordinate   *nuwav1.Coordinate   `json:"coordinate,omitempty"`
	NodeAffinity *corev1.NodeAffinity `json:"nodeAffinity,omitempty"`
	NodeList     *corev1.NodeList     `json:"nodeList,omitempty"`
	NodeNameList []string             `json:"nodeNameList,omitempty"`
}

func newLocalCoordinate(index int, client client.Client, coordinate nuwav1.Coordinate) (*coordinator, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	lc := &coordinator{
		Index:      index,
		Client:     client,
		Coordinate: &coordinate,
	}
	if err := lc.setCoordinateName(); err != nil {
		return nil, err
	}

	if err := lc.findMatchNode(); err != nil {
		return nil, err
	}

	lc.generateAffinity()

	return lc, nil
}

func (l *coordinator) generateAffinity() {
	l.NodeAffinity = &corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						corev1.NodeSelectorRequirement{
							Key:      nuwaRoomFlag,
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{l.Coordinate.Room},
						},
					},
				},
			},
		},
		PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
			corev1.PreferredSchedulingTerm{
				Weight: 100,
				Preference: corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						corev1.NodeSelectorRequirement{
							Key:      nuwaRoomFlag,
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{l.Coordinate.Room},
						},
						corev1.NodeSelectorRequirement{
							Key:      nuwaCabinetFlag,
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{l.Coordinate.Cabinet},
						},
						corev1.NodeSelectorRequirement{
							Key:      nuwaHostFlag,
							Operator: corev1.NodeSelectorOpIn,
							Values:   l.NodeNameList,
						},
					},
				},
			},
		},
	}
}

func (l *coordinator) findMatchNode() (err error) {
	labels, err := l.generateCoordinateMatchLabels()
	if err != nil {
		return err
	}
	l.NodeList = &corev1.NodeList{}
	if err = l.Client.List(context.TODO(), l.NodeList, labels, l.generateHostMatchLabels()); err != nil {
		return err
	}
	for i := range l.NodeList.Items {
		l.NodeNameList = append(l.NodeNameList, l.NodeList.Items[i].Name)
	}

	return
}

func (l *coordinator) generateCoordinateMatchLabels() (clientMatchLabs client.MatchingLabels, err error) {
	if l.Coordinate.Room == "" {
		return nil, ErrNeedAtLeastRoom
	}
	clientMatchLabs = make(client.MatchingLabels)
	clientMatchLabs[nuwaRoomFlag] = l.Coordinate.Room

	if l.Coordinate.Cabinet != "" {
		clientMatchLabs[nuwaCabinetFlag] = l.Coordinate.Cabinet
	}

	return
}

func (l *coordinator) generateHostMatchLabels() (clientMatchLabs client.MatchingLabels) {
	clientMatchLabs = make(client.MatchingLabels)
	if l.Coordinate.Host != "" {
		clientMatchLabs[nuwaHostFlag] = l.Coordinate.Host
	}

	return
}

func (l *coordinator) setCoordinateName() error {
	res := ""
	if l.Coordinate.Room == "" {
		return ErrNeedAtLeastRoom
	}
	res += l.Coordinate.Room
	res += "-"
	if l.Coordinate.Cabinet != "" {
		res += l.Coordinate.Cabinet
	} else {
		res += "Non"
	}

	l.Name = res

	return nil
}

func makeLocalCoordinates(client client.Client, coordinates nuwav1.Coordinates) (cds []*coordinator, err error) {
	for i := range coordinates {
		c, err := newLocalCoordinate(i, client, coordinates[i])
		if err != nil {
			return nil, err
		}
		cds = append(cds, c)
	}

	return
}
