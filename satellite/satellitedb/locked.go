// Code generated by lockedgen using 'go generate'. DO NOT EDIT.

// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellitedb

import (
	"context"
	"crypto"
	"sync"
	"time"

	"github.com/skyrings/skyring-common/tools/uuid"

	"storj.io/storj/pkg/accounting"
	"storj.io/storj/pkg/bwagreement"
	"storj.io/storj/pkg/certdb"
	"storj.io/storj/pkg/datarepair/irreparable"
	"storj.io/storj/pkg/datarepair/queue"
	"storj.io/storj/pkg/overlay"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/console"
	"storj.io/storj/satellite/orders"
)

// locked implements a locking wrapper around satellite.DB.
type locked struct {
	sync.Locker
	db satellite.DB
}

// newLocked returns database wrapped with locker.
func newLocked(db satellite.DB) satellite.DB {
	return &locked{&sync.Mutex{}, db}
}

// BandwidthAgreement returns database for storing bandwidth agreements
func (m *locked) BandwidthAgreement() bwagreement.DB {
	m.Lock()
	defer m.Unlock()
	return &lockedBandwidthAgreement{m.Locker, m.db.BandwidthAgreement()}
}

// lockedBandwidthAgreement implements locking wrapper for bwagreement.DB
type lockedBandwidthAgreement struct {
	sync.Locker
	db bwagreement.DB
}

// DeleteExpired deletes orders that are expired and were created before some time
func (m *lockedBandwidthAgreement) DeleteExpired(ctx context.Context, a1 time.Time, a2 time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.DeleteExpired(ctx, a1, a2)
}

// GetExpired gets orders that are expired and were created before some time
func (m *lockedBandwidthAgreement) GetExpired(ctx context.Context, a1 time.Time, a2 time.Time) ([]bwagreement.SavedOrder, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetExpired(ctx, a1, a2)
}

// GetTotalsSince returns the sum of each bandwidth type after (exluding) a given date range
func (m *lockedBandwidthAgreement) GetTotals(ctx context.Context, a1 time.Time, a2 time.Time) (map[storj.NodeID][]int64, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetTotals(ctx, a1, a2)
}

// GetTotals returns stats about an uplink
func (m *lockedBandwidthAgreement) GetUplinkStats(ctx context.Context, a1 time.Time, a2 time.Time) ([]bwagreement.UplinkStat, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetUplinkStats(ctx, a1, a2)
}

// SaveOrder saves an order for accounting
func (m *lockedBandwidthAgreement) SaveOrder(ctx context.Context, a1 *pb.RenterBandwidthAllocation) error {
	m.Lock()
	defer m.Unlock()
	return m.db.SaveOrder(ctx, a1)
}

// CertDB returns database for storing uplink's public key & ID
func (m *locked) CertDB() certdb.DB {
	m.Lock()
	defer m.Unlock()
	return &lockedCertDB{m.Locker, m.db.CertDB()}
}

// lockedCertDB implements locking wrapper for certdb.DB
type lockedCertDB struct {
	sync.Locker
	db certdb.DB
}

// GetPublicKey gets the public key of uplink corresponding to uplink id
func (m *lockedCertDB) GetPublicKey(ctx context.Context, a1 storj.NodeID) (crypto.PublicKey, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetPublicKey(ctx, a1)
}

// SavePublicKey adds a new bandwidth agreement.
func (m *lockedCertDB) SavePublicKey(ctx context.Context, a1 storj.NodeID, a2 crypto.PublicKey) error {
	m.Lock()
	defer m.Unlock()
	return m.db.SavePublicKey(ctx, a1, a2)
}

// Close closes the database
func (m *locked) Close() error {
	m.Lock()
	defer m.Unlock()
	return m.db.Close()
}

// Console returns database for satellite console
func (m *locked) Console() console.DB {
	m.Lock()
	defer m.Unlock()
	return &lockedConsole{m.Locker, m.db.Console()}
}

// lockedConsole implements locking wrapper for console.DB
type lockedConsole struct {
	sync.Locker
	db console.DB
}

// APIKeys is a getter for APIKeys repository
func (m *lockedConsole) APIKeys() console.APIKeys {
	m.Lock()
	defer m.Unlock()
	return &lockedAPIKeys{m.Locker, m.db.APIKeys()}
}

// lockedAPIKeys implements locking wrapper for console.APIKeys
type lockedAPIKeys struct {
	sync.Locker
	db console.APIKeys
}

// Create creates and stores new APIKeyInfo
func (m *lockedAPIKeys) Create(ctx context.Context, key console.APIKey, info console.APIKeyInfo) (*console.APIKeyInfo, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Create(ctx, key, info)
}

// Delete deletes APIKeyInfo from store
func (m *lockedAPIKeys) Delete(ctx context.Context, id uuid.UUID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, id)
}

// Get retrieves APIKeyInfo with given ID
func (m *lockedAPIKeys) Get(ctx context.Context, id uuid.UUID) (*console.APIKeyInfo, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Get(ctx, id)
}

// GetByKey retrieves APIKeyInfo for given key
func (m *lockedAPIKeys) GetByKey(ctx context.Context, key console.APIKey) (*console.APIKeyInfo, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByKey(ctx, key)
}

// GetByProjectID retrieves list of APIKeys for given projectID
func (m *lockedAPIKeys) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]console.APIKeyInfo, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByProjectID(ctx, projectID)
}

// Update updates APIKeyInfo in store
func (m *lockedAPIKeys) Update(ctx context.Context, key console.APIKeyInfo) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Update(ctx, key)
}

// BucketUsage is a getter for accounting.BucketUsage repository
func (m *lockedConsole) BucketUsage() accounting.BucketUsage {
	m.Lock()
	defer m.Unlock()
	return &lockedBucketUsage{m.Locker, m.db.BucketUsage()}
}

// lockedBucketUsage implements locking wrapper for accounting.BucketUsage
type lockedBucketUsage struct {
	sync.Locker
	db accounting.BucketUsage
}

func (m *lockedBucketUsage) Create(ctx context.Context, rollup accounting.BucketRollup) (*accounting.BucketRollup, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Create(ctx, rollup)
}

func (m *lockedBucketUsage) Delete(ctx context.Context, id uuid.UUID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, id)
}

func (m *lockedBucketUsage) Get(ctx context.Context, id uuid.UUID) (*accounting.BucketRollup, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Get(ctx, id)
}

func (m *lockedBucketUsage) GetPaged(ctx context.Context, cursor *accounting.BucketRollupCursor) ([]accounting.BucketRollup, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetPaged(ctx, cursor)
}

// ProjectMembers is a getter for ProjectMembers repository
func (m *lockedConsole) ProjectMembers() console.ProjectMembers {
	m.Lock()
	defer m.Unlock()
	return &lockedProjectMembers{m.Locker, m.db.ProjectMembers()}
}

// lockedProjectMembers implements locking wrapper for console.ProjectMembers
type lockedProjectMembers struct {
	sync.Locker
	db console.ProjectMembers
}

// Delete is a method for deleting project member by memberID and projectID from the database.
func (m *lockedProjectMembers) Delete(ctx context.Context, memberID uuid.UUID, projectID uuid.UUID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, memberID, projectID)
}

// GetByMemberID is a method for querying project members from the database by memberID.
func (m *lockedProjectMembers) GetByMemberID(ctx context.Context, memberID uuid.UUID) ([]console.ProjectMember, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByMemberID(ctx, memberID)
}

// GetByProjectID is a method for querying project members from the database by projectID, offset and limit.
func (m *lockedProjectMembers) GetByProjectID(ctx context.Context, projectID uuid.UUID, pagination console.Pagination) ([]console.ProjectMember, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByProjectID(ctx, projectID, pagination)
}

// Insert is a method for inserting project member into the database.
func (m *lockedProjectMembers) Insert(ctx context.Context, memberID uuid.UUID, projectID uuid.UUID) (*console.ProjectMember, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Insert(ctx, memberID, projectID)
}

// Projects is a getter for Projects repository
func (m *lockedConsole) Projects() console.Projects {
	m.Lock()
	defer m.Unlock()
	return &lockedProjects{m.Locker, m.db.Projects()}
}

// lockedProjects implements locking wrapper for console.Projects
type lockedProjects struct {
	sync.Locker
	db console.Projects
}

// Delete is a method for deleting project by Id from the database.
func (m *lockedProjects) Delete(ctx context.Context, id uuid.UUID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, id)
}

// Get is a method for querying project from the database by id.
func (m *lockedProjects) Get(ctx context.Context, id uuid.UUID) (*console.Project, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Get(ctx, id)
}

// GetAll is a method for querying all projects from the database.
func (m *lockedProjects) GetAll(ctx context.Context) ([]console.Project, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetAll(ctx)
}

// GetByUserID is a method for querying all projects from the database by userID.
func (m *lockedProjects) GetByUserID(ctx context.Context, userID uuid.UUID) ([]console.Project, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByUserID(ctx, userID)
}

// Insert is a method for inserting project into the database.
func (m *lockedProjects) Insert(ctx context.Context, project *console.Project) (*console.Project, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Insert(ctx, project)
}

// Update is a method for updating project entity.
func (m *lockedProjects) Update(ctx context.Context, project *console.Project) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Update(ctx, project)
}

// RegistrationTokens is a getter for RegistrationTokens repository
func (m *lockedConsole) RegistrationTokens() console.RegistrationTokens {
	m.Lock()
	defer m.Unlock()
	return &lockedRegistrationTokens{m.Locker, m.db.RegistrationTokens()}
}

// lockedRegistrationTokens implements locking wrapper for console.RegistrationTokens
type lockedRegistrationTokens struct {
	sync.Locker
	db console.RegistrationTokens
}

// Create creates new registration token
func (m *lockedRegistrationTokens) Create(ctx context.Context, projectLimit int) (*console.RegistrationToken, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Create(ctx, projectLimit)
}

// GetByOwnerID retrieves RegTokenInfo by ownerID
func (m *lockedRegistrationTokens) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*console.RegistrationToken, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByOwnerID(ctx, ownerID)
}

// GetBySecret retrieves RegTokenInfo with given Secret
func (m *lockedRegistrationTokens) GetBySecret(ctx context.Context, secret console.RegistrationSecret) (*console.RegistrationToken, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetBySecret(ctx, secret)
}

// UpdateOwner updates registration token's owner
func (m *lockedRegistrationTokens) UpdateOwner(ctx context.Context, secret console.RegistrationSecret, ownerID uuid.UUID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateOwner(ctx, secret, ownerID)
}

// ResetPasswordTokens is a getter for ResetPasswordTokens repository
func (m *lockedConsole) ResetPasswordTokens() console.ResetPasswordTokens {
	m.Lock()
	defer m.Unlock()
	return &lockedResetPasswordTokens{m.Locker, m.db.ResetPasswordTokens()}
}

// lockedResetPasswordTokens implements locking wrapper for console.ResetPasswordTokens
type lockedResetPasswordTokens struct {
	sync.Locker
	db console.ResetPasswordTokens
}

// Create creates new reset password token
func (m *lockedResetPasswordTokens) Create(ctx context.Context, ownerID uuid.UUID) (*console.ResetPasswordToken, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Create(ctx, ownerID)
}

// Delete deletes ResetPasswordToken by ResetPasswordSecret
func (m *lockedResetPasswordTokens) Delete(ctx context.Context, secret console.ResetPasswordSecret) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, secret)
}

// GetByOwnerID retrieves ResetPasswordToken by ownerID
func (m *lockedResetPasswordTokens) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*console.ResetPasswordToken, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByOwnerID(ctx, ownerID)
}

// GetBySecret retrieves ResetPasswordToken with given secret
func (m *lockedResetPasswordTokens) GetBySecret(ctx context.Context, secret console.ResetPasswordSecret) (*console.ResetPasswordToken, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetBySecret(ctx, secret)
}

// UsageRollups is a getter for UsageRollups repository
func (m *lockedConsole) UsageRollups() console.UsageRollups {
	m.Lock()
	defer m.Unlock()
	return &lockedUsageRollups{m.Locker, m.db.UsageRollups()}
}

// lockedUsageRollups implements locking wrapper for console.UsageRollups
type lockedUsageRollups struct {
	sync.Locker
	db console.UsageRollups
}

func (m *lockedUsageRollups) GetBucketTotals(ctx context.Context, projectID uuid.UUID, cursor console.BucketUsageCursor, since time.Time, before time.Time) (*console.BucketUsagePage, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetBucketTotals(ctx, projectID, cursor, since, before)
}

func (m *lockedUsageRollups) GetBucketUsageRollups(ctx context.Context, projectID uuid.UUID, since time.Time, before time.Time) ([]console.BucketUsageRollup, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetBucketUsageRollups(ctx, projectID, since, before)
}

func (m *lockedUsageRollups) GetProjectTotal(ctx context.Context, projectID uuid.UUID, since time.Time, before time.Time) (*console.ProjectUsage, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetProjectTotal(ctx, projectID, since, before)
}

// Users is a getter for Users repository
func (m *lockedConsole) Users() console.Users {
	m.Lock()
	defer m.Unlock()
	return &lockedUsers{m.Locker, m.db.Users()}
}

// lockedUsers implements locking wrapper for console.Users
type lockedUsers struct {
	sync.Locker
	db console.Users
}

// Delete is a method for deleting user by Id from the database.
func (m *lockedUsers) Delete(ctx context.Context, id uuid.UUID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, id)
}

// Get is a method for querying user from the database by id.
func (m *lockedUsers) Get(ctx context.Context, id uuid.UUID) (*console.User, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Get(ctx, id)
}

// GetByEmail is a method for querying user by email from the database.
func (m *lockedUsers) GetByEmail(ctx context.Context, email string) (*console.User, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetByEmail(ctx, email)
}

// Insert is a method for inserting user into the database.
func (m *lockedUsers) Insert(ctx context.Context, user *console.User) (*console.User, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Insert(ctx, user)
}

// Update is a method for updating user entity.
func (m *lockedUsers) Update(ctx context.Context, user *console.User) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Update(ctx, user)
}

// CreateSchema sets the schema
func (m *locked) CreateSchema(schema string) error {
	m.Lock()
	defer m.Unlock()
	return m.db.CreateSchema(schema)
}

// CreateTables initializes the database
func (m *locked) CreateTables() error {
	m.Lock()
	defer m.Unlock()
	return m.db.CreateTables()
}

// DropSchema drops the schema
func (m *locked) DropSchema(schema string) error {
	m.Lock()
	defer m.Unlock()
	return m.db.DropSchema(schema)
}

// Irreparable returns database for failed repairs
func (m *locked) Irreparable() irreparable.DB {
	m.Lock()
	defer m.Unlock()
	return &lockedIrreparable{m.Locker, m.db.Irreparable()}
}

// lockedIrreparable implements locking wrapper for irreparable.DB
type lockedIrreparable struct {
	sync.Locker
	db irreparable.DB
}

// Delete removes irreparable segment info based on segmentPath.
func (m *lockedIrreparable) Delete(ctx context.Context, segmentPath []byte) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, segmentPath)
}

// Get returns irreparable segment info based on segmentPath.
func (m *lockedIrreparable) Get(ctx context.Context, segmentPath []byte) (*pb.IrreparableSegment, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Get(ctx, segmentPath)
}

// GetLimited number of segments from offset
func (m *lockedIrreparable) GetLimited(ctx context.Context, limit int, offset int64) ([]*pb.IrreparableSegment, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetLimited(ctx, limit, offset)
}

// IncrementRepairAttempts increments the repair attempts.
func (m *lockedIrreparable) IncrementRepairAttempts(ctx context.Context, segmentInfo *pb.IrreparableSegment) error {
	m.Lock()
	defer m.Unlock()
	return m.db.IncrementRepairAttempts(ctx, segmentInfo)
}

// Orders returns database for orders
func (m *locked) Orders() orders.DB {
	m.Lock()
	defer m.Unlock()
	return &lockedOrders{m.Locker, m.db.Orders()}
}

// lockedOrders implements locking wrapper for orders.DB
type lockedOrders struct {
	sync.Locker
	db orders.DB
}

// CreateSerialInfo creates serial number entry in database
func (m *lockedOrders) CreateSerialInfo(ctx context.Context, serialNumber storj.SerialNumber, bucketID []byte, limitExpiration time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.CreateSerialInfo(ctx, serialNumber, bucketID, limitExpiration)
}

// GetBucketBandwidth gets total bucket bandwidth from period of time
func (m *lockedOrders) GetBucketBandwidth(ctx context.Context, bucketID []byte, from time.Time, to time.Time) (int64, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetBucketBandwidth(ctx, bucketID, from, to)
}

// GetStorageNodeBandwidth gets total storage node bandwidth from period of time
func (m *lockedOrders) GetStorageNodeBandwidth(ctx context.Context, nodeID storj.NodeID, from time.Time, to time.Time) (int64, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetStorageNodeBandwidth(ctx, nodeID, from, to)
}

// UnuseSerialNumber removes pair serial number -> storage node id from database
func (m *lockedOrders) UnuseSerialNumber(ctx context.Context, serialNumber storj.SerialNumber, storageNodeID storj.NodeID) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UnuseSerialNumber(ctx, serialNumber, storageNodeID)
}

// UpdateBucketBandwidthAllocation updates 'allocated' bandwidth for given bucket
func (m *lockedOrders) UpdateBucketBandwidthAllocation(ctx context.Context, bucketID []byte, action pb.PieceAction, amount int64, intervalStart time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateBucketBandwidthAllocation(ctx, bucketID, action, amount, intervalStart)
}

// UpdateBucketBandwidthInline updates 'inline' bandwidth for given bucket
func (m *lockedOrders) UpdateBucketBandwidthInline(ctx context.Context, bucketID []byte, action pb.PieceAction, amount int64, intervalStart time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateBucketBandwidthInline(ctx, bucketID, action, amount, intervalStart)
}

// UpdateBucketBandwidthSettle updates 'settled' bandwidth for given bucket
func (m *lockedOrders) UpdateBucketBandwidthSettle(ctx context.Context, bucketID []byte, action pb.PieceAction, amount int64, intervalStart time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateBucketBandwidthSettle(ctx, bucketID, action, amount, intervalStart)
}

// UpdateStoragenodeBandwidthAllocation updates 'allocated' bandwidth for given storage node
func (m *lockedOrders) UpdateStoragenodeBandwidthAllocation(ctx context.Context, storageNode storj.NodeID, action pb.PieceAction, amount int64, intervalStart time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateStoragenodeBandwidthAllocation(ctx, storageNode, action, amount, intervalStart)
}

// UpdateStoragenodeBandwidthSettle updates 'settled' bandwidth for given storage node
func (m *lockedOrders) UpdateStoragenodeBandwidthSettle(ctx context.Context, storageNode storj.NodeID, action pb.PieceAction, amount int64, intervalStart time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateStoragenodeBandwidthSettle(ctx, storageNode, action, amount, intervalStart)
}

// UseSerialNumber creates serial number entry in database
func (m *lockedOrders) UseSerialNumber(ctx context.Context, serialNumber storj.SerialNumber, storageNodeID storj.NodeID) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.UseSerialNumber(ctx, serialNumber, storageNodeID)
}

// OverlayCache returns database for caching overlay information
func (m *locked) OverlayCache() overlay.DB {
	m.Lock()
	defer m.Unlock()
	return &lockedOverlayCache{m.Locker, m.db.OverlayCache()}
}

// lockedOverlayCache implements locking wrapper for overlay.DB
type lockedOverlayCache struct {
	sync.Locker
	db overlay.DB
}

// CreateStats initializes the stats for node.
func (m *lockedOverlayCache) CreateStats(ctx context.Context, nodeID storj.NodeID, initial *overlay.NodeStats) (stats *overlay.NodeStats, err error) {
	m.Lock()
	defer m.Unlock()
	return m.db.CreateStats(ctx, nodeID, initial)
}

// Get looks up the node by nodeID
func (m *lockedOverlayCache) Get(ctx context.Context, nodeID storj.NodeID) (*overlay.NodeDossier, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Get(ctx, nodeID)
}

// AllUnreliableOrOffline returns all unreliable or offlines node, independent of new
func (m *lockedOverlayCache) AllUnreliableOrOffline(ctx context.Context, criteria *overlay.NodeCriteria) (badNodes map[storj.NodeID]struct{}, err error) {
	m.Lock()
	defer m.Unlock()
	return m.db.AllUnreliableOrOffline(ctx, criteria)
}

// UnreliableOrOffline filters a set of nodes to unhealth or offlines node, independent of new
func (m *lockedOverlayCache) UnreliableOrOffline(ctx context.Context, a1 *overlay.NodeCriteria, a2 storj.NodeIDList) (storj.NodeIDList, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.UnreliableOrOffline(ctx, a1, a2)
}

// Paginate will page through the database nodes
func (m *lockedOverlayCache) Paginate(ctx context.Context, offset int64, limit int) ([]*overlay.NodeDossier, bool, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Paginate(ctx, offset, limit)
}

// SelectNewStorageNodes looks up nodes based on new node criteria
func (m *lockedOverlayCache) SelectNewStorageNodes(ctx context.Context, count int, criteria *overlay.NodeCriteria) ([]*pb.Node, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.SelectNewStorageNodes(ctx, count, criteria)
}

// SelectStorageNodes looks up nodes based on criteria
func (m *lockedOverlayCache) SelectStorageNodes(ctx context.Context, count int, criteria *overlay.NodeCriteria) ([]*pb.Node, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.SelectStorageNodes(ctx, count, criteria)
}

// Update updates node address
func (m *lockedOverlayCache) UpdateAddress(ctx context.Context, value *pb.Node) error {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateAddress(ctx, value)
}

// UpdateNodeInfo updates node dossier with info requested from the node itself like node type, email, wallet, capacity, and version.
func (m *lockedOverlayCache) UpdateNodeInfo(ctx context.Context, node storj.NodeID, nodeInfo *pb.InfoResponse) (stats *overlay.NodeDossier, err error) {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateNodeInfo(ctx, node, nodeInfo)
}

// UpdateStats all parts of single storagenode's stats.
func (m *lockedOverlayCache) UpdateStats(ctx context.Context, request *overlay.UpdateRequest) (stats *overlay.NodeStats, err error) {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateStats(ctx, request)
}

// UpdateUptime updates a single storagenode's uptime stats.
func (m *lockedOverlayCache) UpdateUptime(ctx context.Context, nodeID storj.NodeID, isUp bool) (stats *overlay.NodeStats, err error) {
	m.Lock()
	defer m.Unlock()
	return m.db.UpdateUptime(ctx, nodeID, isUp)
}

// ProjectAccounting returns database for storing information about project data use
func (m *locked) ProjectAccounting() accounting.ProjectAccounting {
	m.Lock()
	defer m.Unlock()
	return &lockedProjectAccounting{m.Locker, m.db.ProjectAccounting()}
}

// lockedProjectAccounting implements locking wrapper for accounting.ProjectAccounting
type lockedProjectAccounting struct {
	sync.Locker
	db accounting.ProjectAccounting
}

// CreateStorageTally creates a record for BucketStorageTally in the accounting DB table
func (m *lockedProjectAccounting) CreateStorageTally(ctx context.Context, tally accounting.BucketStorageTally) error {
	m.Lock()
	defer m.Unlock()
	return m.db.CreateStorageTally(ctx, tally)
}

// GetAllocatedBandwidthTotal returns the sum of GET bandwidth usage allocated for a projectID in the past time frame
func (m *lockedProjectAccounting) GetAllocatedBandwidthTotal(ctx context.Context, bucketID []byte, from time.Time) (int64, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetAllocatedBandwidthTotal(ctx, bucketID, from)
}

// GetStorageTotals returns the current inline and remote storage usage for a projectID
func (m *lockedProjectAccounting) GetStorageTotals(ctx context.Context, projectID uuid.UUID) (int64, int64, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetStorageTotals(ctx, projectID)
}

// SaveTallies saves the latest project info
func (m *lockedProjectAccounting) SaveTallies(ctx context.Context, intervalStart time.Time, bucketTallies map[string]*accounting.BucketTally) ([]accounting.BucketTally, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.SaveTallies(ctx, intervalStart, bucketTallies)
}

// RepairQueue returns queue for segments that need repairing
func (m *locked) RepairQueue() queue.RepairQueue {
	m.Lock()
	defer m.Unlock()
	return &lockedRepairQueue{m.Locker, m.db.RepairQueue()}
}

// lockedRepairQueue implements locking wrapper for queue.RepairQueue
type lockedRepairQueue struct {
	sync.Locker
	db queue.RepairQueue
}

// Delete removes an injured segment.
func (m *lockedRepairQueue) Delete(ctx context.Context, s *pb.InjuredSegment) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Delete(ctx, s)
}

// Insert adds an injured segment.
func (m *lockedRepairQueue) Insert(ctx context.Context, s *pb.InjuredSegment) error {
	m.Lock()
	defer m.Unlock()
	return m.db.Insert(ctx, s)
}

// Select gets an injured segment.
func (m *lockedRepairQueue) Select(ctx context.Context) (*pb.InjuredSegment, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.Select(ctx)
}

// SelectN lists limit amount of injured segments.
func (m *lockedRepairQueue) SelectN(ctx context.Context, limit int) ([]pb.InjuredSegment, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.SelectN(ctx, limit)
}

// StoragenodeAccounting returns database for storing information about storagenode use
func (m *locked) StoragenodeAccounting() accounting.StoragenodeAccounting {
	m.Lock()
	defer m.Unlock()
	return &lockedStoragenodeAccounting{m.Locker, m.db.StoragenodeAccounting()}
}

// lockedStoragenodeAccounting implements locking wrapper for accounting.StoragenodeAccounting
type lockedStoragenodeAccounting struct {
	sync.Locker
	db accounting.StoragenodeAccounting
}

// DeleteTalliesBefore deletes all tallies prior to some time
func (m *lockedStoragenodeAccounting) DeleteTalliesBefore(ctx context.Context, latestRollup time.Time) error {
	m.Lock()
	defer m.Unlock()
	return m.db.DeleteTalliesBefore(ctx, latestRollup)
}

// GetBandwidthSince retrieves all bandwidth rollup entires since latestRollup
func (m *lockedStoragenodeAccounting) GetBandwidthSince(ctx context.Context, latestRollup time.Time) ([]*accounting.StoragenodeBandwidthRollup, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetBandwidthSince(ctx, latestRollup)
}

// GetTallies retrieves all tallies
func (m *lockedStoragenodeAccounting) GetTallies(ctx context.Context) ([]*accounting.StoragenodeStorageTally, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetTallies(ctx)
}

// GetTalliesSince retrieves all tallies since latestRollup
func (m *lockedStoragenodeAccounting) GetTalliesSince(ctx context.Context, latestRollup time.Time) ([]*accounting.StoragenodeStorageTally, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.GetTalliesSince(ctx, latestRollup)
}

// LastTimestamp records and returns the latest last tallied time.
func (m *lockedStoragenodeAccounting) LastTimestamp(ctx context.Context, timestampType string) (time.Time, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.LastTimestamp(ctx, timestampType)
}

// QueryPaymentInfo queries Nodes and Accounting_Rollup on nodeID
func (m *lockedStoragenodeAccounting) QueryPaymentInfo(ctx context.Context, start time.Time, end time.Time) ([]*accounting.CSVRow, error) {
	m.Lock()
	defer m.Unlock()
	return m.db.QueryPaymentInfo(ctx, start, end)
}

// SaveRollup records tally and bandwidth rollup aggregations to the database
func (m *lockedStoragenodeAccounting) SaveRollup(ctx context.Context, latestTally time.Time, stats accounting.RollupStats) error {
	m.Lock()
	defer m.Unlock()
	return m.db.SaveRollup(ctx, latestTally, stats)
}

// SaveTallies records tallies of data at rest
func (m *lockedStoragenodeAccounting) SaveTallies(ctx context.Context, latestTally time.Time, nodeData map[storj.NodeID]float64) error {
	m.Lock()
	defer m.Unlock()
	return m.db.SaveTallies(ctx, latestTally, nodeData)
}
