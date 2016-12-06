package client

import (
	artifactory "github.com/appscode/api/artifactory/v1beta1"
	auth "github.com/appscode/api/auth/v1beta1"
	backup "github.com/appscode/api/backup/v1beta1"
	billing "github.com/appscode/api/billing/v1beta1"
	ca "github.com/appscode/api/certificate/v1beta1"
	ci "github.com/appscode/api/ci/v1beta1"
	db "github.com/appscode/api/db/v1beta1"
	glusterfs "github.com/appscode/api/glusterfs/v1beta1"
	kubernetes "github.com/appscode/api/kubernetes/v1beta1"
	pv "github.com/appscode/api/pv/v1beta1"
	"google.golang.org/grpc"
)

// multi client services are grouped by there main client. the api service
// clients are wrapped around with sub-service.
type multiClientInterface interface {
	Artifactory() *artifactoryService
	Authentication() *authenticationService
	Backup() *backupService
	Billing() *billingService
	CA() *caService
	CI() *ciService
	GlusterFS() *glusterFSService
	Kubernetes() *kubernetesService
	PV() *pvService
	DB() *dbService
}

type multiClientServices struct {
	artifactoryClient    *artifactoryService
	authenticationClient *authenticationService
	backupClient         *backupService
	billingClient        *billingService
	caClient             *caService
	ciClient             *ciService
	glusterfsClient      *glusterFSService
	kubernetesClient     *kubernetesService
	pvClient             *pvService
	dbClient             *dbService
}

func newMultiClientService(conn *grpc.ClientConn) multiClientInterface {
	return &multiClientServices{
		artifactoryClient: &artifactoryService{
			artifactClient: artifactory.NewArtifactsClient(conn),
			versionClient:  artifactory.NewVersionsClient(conn),
		},
		authenticationClient: &authenticationService{
			authenticationClient: auth.NewAuthenticationClient(conn),
			conduitClient:        auth.NewConduitClient(conn),
		},
		backupClient: &backupService{
			backupServerClient: backup.NewServersClient(conn),
			backupClientClient: backup.NewClientsClient(conn),
		},
		billingClient: &billingService{
			paymentMethodClient: billing.NewPaymentMethodsClient(conn),
			subscriptionClient:  billing.NewSubscriptionsClient(conn),
			purchaseClient:      billing.NewPurchasesClient(conn),
			chargeClient:        billing.NewChargesClient(conn),
			quotaClient:         billing.NewQuotasClient(conn),
		},
		caClient: &caService{
			certificateClient: ca.NewCertificatesClient(conn),
		},
		ciClient: &ciService{
			buildClient:  ci.NewBuildsClient(conn),
			jobClient:    ci.NewJobsClient(conn),
			masterClient: ci.NewMasterClient(conn),
			agentClient:  ci.NewAgentsClient(conn),
		},
		glusterfsClient: &glusterFSService{
			clusterClient: glusterfs.NewClustersClient(conn),
			volumeClient:  glusterfs.NewVolumesClient(conn),
		},
		kubernetesClient: &kubernetesService{
			kubernetesClient: kubernetes.NewClientsClient(conn),
			clusterClient:    kubernetes.NewClustersClient(conn),
			eventsClient:     kubernetes.NewEventsClient(conn),
			metdataClient:    kubernetes.NewMetadataClient(conn),
			incidentClient:   kubernetes.NewIncidentsClient(conn),
		},
		pvClient: &pvService{
			disk: pv.NewDisksClient(conn),
			pv:   pv.NewPersistentVolumesClient(conn),
			pvc:  pv.NewPersistentVolumeClaimsClient(conn),
		},
		dbClient: &dbService{
			database: db.NewDatabasesClient(conn),
			snapshot: db.NewSnapshotsClient(conn),
		},
	}
}

func (s *multiClientServices) Artifactory() *artifactoryService {
	return s.artifactoryClient
}

func (s *multiClientServices) Authentication() *authenticationService {
	return s.authenticationClient
}

func (s *multiClientServices) Backup() *backupService {
	return s.backupClient
}

func (s *multiClientServices) Billing() *billingService {
	return s.billingClient
}

func (s *multiClientServices) CA() *caService {
	return s.caClient
}

func (s *multiClientServices) CI() *ciService {
	return s.ciClient
}

func (s *multiClientServices) GlusterFS() *glusterFSService {
	return s.glusterfsClient
}

func (s *multiClientServices) Kubernetes() *kubernetesService {
	return s.kubernetesClient
}

func (s *multiClientServices) PV() *pvService {
	return s.pvClient
}

func (s *multiClientServices) DB() *dbService {
	return s.dbClient
}

// original service clients that needs to exposed under grouped wrapper
// services.
type artifactoryService struct {
	artifactClient artifactory.ArtifactsClient
	versionClient  artifactory.VersionsClient
}

func (a *artifactoryService) Artifacts() artifactory.ArtifactsClient {
	return a.artifactClient
}

func (a *artifactoryService) Versions() artifactory.VersionsClient {
	return a.versionClient
}

type authenticationService struct {
	authenticationClient auth.AuthenticationClient
	conduitClient        auth.ConduitClient
}

func (a *authenticationService) Authentication() auth.AuthenticationClient {
	return a.authenticationClient
}

func (a *authenticationService) Conduit() auth.ConduitClient {
	return a.conduitClient
}

type backupService struct {
	backupServerClient backup.ServersClient
	backupClientClient backup.ClientsClient
}

func (b *backupService) Server() backup.ServersClient {
	return b.backupServerClient
}

func (b *backupService) Client() backup.ClientsClient {
	return b.backupClientClient
}

type billingService struct {
	paymentMethodClient billing.PaymentMethodsClient
	subscriptionClient  billing.SubscriptionsClient
	purchaseClient      billing.PurchasesClient
	chargeClient        billing.ChargesClient
	quotaClient         billing.QuotasClient
}

func (b *billingService) Charge() billing.ChargesClient {
	return b.chargeClient
}

func (b *billingService) Subscription() billing.SubscriptionsClient {
	return b.subscriptionClient
}

func (b *billingService) Quota() billing.QuotasClient {
	return b.quotaClient
}

type caService struct {
	certificateClient ca.CertificatesClient
}

func (c *caService) CertificatesClient() ca.CertificatesClient {
	return c.certificateClient
}

type ciService struct {
	buildClient  ci.BuildsClient
	jobClient    ci.JobsClient
	masterClient ci.MasterClient
	agentClient  ci.AgentsClient
}

func (c *ciService) Build() ci.BuildsClient {
	return c.buildClient
}

func (c *ciService) Job() ci.JobsClient {
	return c.jobClient
}

func (c *ciService) Master() ci.MasterClient {
	return c.masterClient
}

func (c *ciService) Agent() ci.AgentsClient {
	return c.agentClient
}

type glusterFSService struct {
	clusterClient glusterfs.ClustersClient
	volumeClient  glusterfs.VolumesClient
}

func (g *glusterFSService) Cluster() glusterfs.ClustersClient {
	return g.clusterClient
}

func (g *glusterFSService) Volume() glusterfs.VolumesClient {
	return g.volumeClient
}

type kubernetesService struct {
	kubernetesClient kubernetes.ClientsClient
	clusterClient    kubernetes.ClustersClient
	eventsClient     kubernetes.EventsClient
	metdataClient    kubernetes.MetadataClient
	incidentClient   kubernetes.IncidentsClient
}

func (k *kubernetesService) Client() kubernetes.ClientsClient {
	return k.kubernetesClient
}

func (k *kubernetesService) Cluster() kubernetes.ClustersClient {
	return k.clusterClient
}

func (k *kubernetesService) Event() kubernetes.EventsClient {
	return k.eventsClient
}

func (k *kubernetesService) Metadata() kubernetes.MetadataClient {
	return k.metdataClient
}

func (a *kubernetesService) Incident() kubernetes.IncidentsClient {
	return a.incidentClient
}

type pvService struct {
	disk pv.DisksClient
	pv   pv.PersistentVolumesClient
	pvc  pv.PersistentVolumeClaimsClient
}

func (p *pvService) Disk() pv.DisksClient {
	return p.disk
}

func (p *pvService) PersistentVolume() pv.PersistentVolumesClient {
	return p.pv
}

func (p *pvService) PersistentVolumeClaim() pv.PersistentVolumeClaimsClient {
	return p.pvc
}

type dbService struct {
	database db.DatabasesClient
	snapshot db.SnapshotsClient
}

func (p *dbService) Database() db.DatabasesClient {
	return p.database
}

func (p *dbService) Snapshot() db.SnapshotsClient {
	return p.snapshot
}
